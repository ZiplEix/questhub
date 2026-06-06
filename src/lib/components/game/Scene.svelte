<script lang="ts">
    import { activeBoardStore, sendPing, onPingReceived } from "$lib/websocket";
    import { Compass, Minus, Plus } from "lucide-svelte";
    import { fade } from "svelte/transition";
    import { onMount, untrack } from "svelte";
    import { page } from "$app/state";
    import { fetchPlayers, fetchGame } from "$lib/api";
    import { authClient } from "$lib/auth-client";


    let {
        isGM = false
    } = $props<{
        isGM?: boolean;
    }>();

    // Derived active map url from store
    let activeMapUrl = $derived($activeBoardStore?.image_url || "");

    // Pings State
    let pings = $state<{ x: number; y: number; id: number; color?: string }[]>([]);
    let myColor = "#E07A5F"; // default ping color


    // Zoom and Pan States
    let zoom = $state(1);
    let pan = { x: 0, y: 0 };

    let isPanning = $state(false);
    let startX = 0;
    let startY = 0;
    let startPanX = 0;
    let startPanY = 0;

    // Canvas size bindings
    let width = $state(0);
    let height = $state(0);
    let pixelRatio = $state(1);

    let containerEl = $state<HTMLElement | null>(null);
    let canvasEl = $state<HTMLCanvasElement | null>(null);

    // Image loading (initialized to undefined for SSR safety)
    let img: HTMLImageElement | undefined = undefined;
    let imgLoaded = $state(false);
    let naturalWidth = 0;
    let naturalHeight = 0;

    let animId: number | null = null;

    // Load image when activeMapUrl changes
    $effect(() => {
        if (activeMapUrl) {
            imgLoaded = false;
            img = new Image(); // Client-side instantiation is safe
            img.src = activeMapUrl;
            img.onload = () => {
                imgLoaded = true;
                if (!img) return;
                naturalWidth = img.naturalWidth;
                naturalHeight = img.naturalHeight;
                
                // Centering map by default and scaling to fit
                if (width && height) {
                    const fitScale = Math.min((width * 0.9) / naturalWidth, (height * 0.9) / naturalHeight);
                    zoom = Math.min(Math.max(fitScale, 0.15), 2);
                    pan.x = width / 2;
                    pan.y = height / 2;
                }
                draw();
            };
        } else {
            imgLoaded = false;
        }
    });

    // Handle canvas resizing and redrawing
    $effect(() => {
        if (width && height && imgLoaded) {
            const _ = pixelRatio; // track pixelRatio
            untrack(() => draw());
        }
    });

    onMount(() => {
        pixelRatio = window.devicePixelRatio || 1;
        if (containerEl) {
            containerEl.addEventListener("wheel", handleWheel, { passive: false });
        }

        // Fetch current user's ping color
        const gameId = page.params.id;
        console.log("Scene: gameId is", gameId);
        if (gameId) {
            authClient.getSession().then(({ data: sessionData }) => {
                console.log("Scene: sessionData is", sessionData);
                if (sessionData?.user) {
                    const userId = sessionData.user.id;
                    console.log("Scene: userId is", userId);
                    fetchGame(gameId).then((gameData) => {
                        console.log("Scene: gameData is", gameData);
                        if (userId === gameData.gm_id) {
                            myColor = gameData.gm_ping_color || "#E07A5F";
                            console.log("Scene: User is GM, myColor set to", myColor);
                        } else {
                            fetchPlayers(gameId).then((playersList) => {
                                console.log("Scene: playersList is", playersList);
                                const me = playersList.find(p => p.user_id === userId);
                                console.log("Scene: matched me is", me);
                                if (me && me.ping_color) {
                                    myColor = me.ping_color;
                                    console.log("Scene: Player myColor set to", myColor);
                                }
                            }).catch(err => console.error("Failed to load players for ping color:", err));
                        }
                    }).catch(err => console.error("Failed to load game for ping color:", err));
                }
            }).catch(err => console.error("Failed to get session for ping color:", err));
        }

        const unsubscribePing = onPingReceived((pingData) => {
            if (!imgLoaded) return;
            if (pings.some(p => p.id === pingData.id)) return;
            
            const x = pingData.x * naturalWidth;
            const y = pingData.y * naturalHeight;
            pings = [...pings, { x, y, id: pingData.id, color: pingData.color }];
            startAnimationLoop();
        });

        return () => {
            unsubscribePing();
            if (containerEl) {
                containerEl.removeEventListener("wheel", handleWheel);
            }
            if (animId !== null) {
                cancelAnimationFrame(animId);
            }
        };
    });

    function handleMouseDown(e: MouseEvent) {
        if (!activeMapUrl) return;

        if (e.button === 1) {
            // Middle-click to pan
            e.preventDefault();
            isPanning = true;
            startX = e.clientX;
            startY = e.clientY;
            startPanX = pan.x;
            startPanY = pan.y;
        } else if (e.button === 0) {
            // Left-click to ping
            handleSceneClick(e);
        }
    }

    function handleMouseMove(e: MouseEvent) {
        if (!isPanning) return;
        const dx = e.clientX - startX;
        const dy = e.clientY - startY;
        pan.x = startPanX + dx;
        pan.y = startPanY + dy;
        
        if (animId === null) {
            draw();
        }
    }

    function handleMouseUp(e: MouseEvent) {
        if (e.button === 1) {
            isPanning = false;
        }
    }

    function handleWheel(e: WheelEvent) {
        if (!activeMapUrl) return;
        e.preventDefault();

        if (e.ctrlKey) {
            // ZOOM (pinch-to-zoom or Ctrl + scroll wheel)
            const zoomIntensity = 0.05;
            const zoomFactor = e.deltaY < 0 ? (1 + zoomIntensity) : (1 - zoomIntensity);
            
            const newZoom = Math.min(Math.max(zoom * zoomFactor, 0.15), 8); // clamp between 15% and 800%
            
            if (newZoom !== zoom && containerEl) {
                const rect = containerEl.getBoundingClientRect();
                const mouseX = e.clientX - rect.left;
                const mouseY = e.clientY - rect.top;
                
                // Adjust pan to zoom relative to mouse pointer position
                pan.x = mouseX - (mouseX - pan.x) * (newZoom / zoom);
                pan.y = mouseY - (mouseY - pan.y) * (newZoom / zoom);
                zoom = newZoom;
            }
        } else {
            // PAN (trackpad scroll or mouse scroll)
            if (e.shiftKey) {
                pan.x -= e.deltaY;
            } else {
                pan.x -= e.deltaX;
                pan.y -= e.deltaY;
            }
        }

        if (animId === null) {
            draw();
        }
    }

    function handleSceneClick(e: MouseEvent) {
        if (!activeMapUrl || !imgLoaded || isPanning || !width || !height) return;

        const halfW = naturalWidth / 2;
        const halfH = naturalHeight / 2;

        if (containerEl) {
            const rect = containerEl.getBoundingClientRect();
            const mouseX = e.clientX - rect.left;
            const mouseY = e.clientY - rect.top;

            // Convert viewport mouse coordinate to image space coordinate
            const x_image_centered = (mouseX - pan.x) / zoom;
            const y_image_centered = (mouseY - pan.y) / zoom;

            const x = x_image_centered + halfW;
            const y = y_image_centered + halfH;

            // Verify click is inside image bounds
            if (x < 0 || x > naturalWidth || y < 0 || y > naturalHeight) {
                return;
            }

            const id = Date.now();
            console.log("Scene: Click detected. Local myColor is:", myColor);
            pings = [...pings, { x, y, id, color: myColor }];
            
            // Broadcast ping using normalized coordinates and custom color
            console.log("Scene: Broadcasting ping with color:", myColor);
            sendPing(x / naturalWidth, y / naturalHeight, myColor);
            
            // Start animation loop for pings
            startAnimationLoop();
        }
    }

    function startAnimationLoop() {
        if (animId !== null) return;
        
        function tick() {
            const now = Date.now();
            pings = pings.filter(p => now - p.id < 2000);
            
            draw();
            
            if (pings.length > 0) {
                animId = requestAnimationFrame(tick);
            } else {
                animId = null;
            }
        }
        
        animId = requestAnimationFrame(tick);
    }

    function draw() {
        if (!canvasEl || !imgLoaded || !width || !height || !img) return;
        const ctx = canvasEl.getContext("2d");
        if (!ctx) return;

        // Clear canvas with pixel ratio scaling
        ctx.clearRect(0, 0, width * pixelRatio, height * pixelRatio);

        // Apply transformations
        ctx.save();
        ctx.scale(pixelRatio, pixelRatio);
        ctx.translate(pan.x, pan.y);
        ctx.scale(zoom, zoom);

        // Draw image centered at (0, 0)
        ctx.drawImage(img, -naturalWidth / 2, -naturalHeight / 2, naturalWidth, naturalHeight);

        // Grid layout spacing
        const spacing = 40;
        const halfW = naturalWidth / 2;
        const halfH = naturalHeight / 2;

        // Draw subtle border around map
        ctx.strokeStyle = "rgba(0, 0, 0, 0.25)";
        ctx.lineWidth = 1 / zoom;
        ctx.strokeRect(-halfW, -halfH, naturalWidth, naturalHeight);

        // Subtle grid dots (double layered for visibility on any background)
        for (let x = -Math.floor(halfW / spacing) * spacing; x <= halfW; x += spacing) {
            for (let y = -Math.floor(halfH / spacing) * spacing; y <= halfH; y += spacing) {
                ctx.fillStyle = "rgba(0, 0, 0, 0.15)";
                ctx.beginPath();
                ctx.arc(x, y, 1.2, 0, Math.PI * 2);
                ctx.fill();
                
                ctx.fillStyle = "rgba(255, 255, 255, 0.25)";
                ctx.beginPath();
                ctx.arc(x + 0.5, y + 0.5, 0.8, 0, Math.PI * 2);
                ctx.fill();
            }
        }

        // Draw active pings
        const now = Date.now();
        for (const ping of pings) {
            const age = (now - ping.id) / 1000;
            if (age >= 2.0) continue;
            
            const px = ping.x - halfW;
            const py = ping.y - halfH;
            
            const opacity = 1 - age / 2.0;
            const radius = age * 40; // grows to 80px
            const pingColor = ping.color || "#E07A5F";
            
            ctx.save();
            ctx.globalAlpha = opacity;
            
            ctx.strokeStyle = pingColor;
            ctx.lineWidth = 3 / zoom;
            ctx.beginPath();
            ctx.arc(px, py, radius, 0, Math.PI * 2);
            ctx.stroke();
            
            ctx.fillStyle = pingColor;
            ctx.beginPath();
            ctx.arc(px, py, 6 / zoom, 0, Math.PI * 2);
            ctx.fill();
            
            ctx.restore();
        }

        ctx.restore();
    }

    function zoomIn() {
        const newZoom = Math.min(zoom * 1.25, 8);
        adjustZoomCenter(newZoom);
    }

    function zoomOut() {
        const newZoom = Math.max(zoom / 1.25, 0.15);
        adjustZoomCenter(newZoom);
    }

    function resetZoom() {
        zoom = 1;
        pan.x = width / 2;
        pan.y = height / 2;
        draw();
    }

    function adjustZoomCenter(newZoom: number) {
        if (newZoom !== zoom && width && height) {
            // Zoom centered relative to canvas center
            const centerX = width / 2;
            const centerY = height / 2;
            
            pan.x = centerX - (centerX - pan.x) * (newZoom / zoom);
            pan.y = centerY - (centerY - pan.y) * (newZoom / zoom);
            zoom = newZoom;
            
            if (animId === null) {
                draw();
            }
        }
    }
</script>

<svelte:window 
    onmousemove={handleMouseMove} 
    onmouseup={handleMouseUp} 
/>

<div
    bind:this={containerEl}
    bind:offsetWidth={width}
    bind:offsetHeight={height}
    class="relative w-full h-full bg-stone-950 overflow-hidden select-none"
    onmousedown={handleMouseDown}
    role="button"
    tabindex="0"
    style="cursor: {isPanning ? 'grabbing' : activeMapUrl ? 'grab' : 'default'};"
>
    {#if activeMapUrl}
        <!-- Canvas Layer -->
        <canvas
            bind:this={canvasEl}
            width={width * pixelRatio}
            height={height * pixelRatio}
            class="absolute inset-0 block"
            style="width: {width}px; height: {height}px;"
        ></canvas>

        <!-- Premium Controls Overlay (Zoom +, Zoom -, Reset) -->
        <div class="absolute bottom-4 right-4 flex gap-1 bg-black/60 backdrop-blur-md p-1.5 rounded-xl border border-white/10 shadow-xl items-center z-30">
            <button
                onclick={zoomOut}
                class="p-1.5 text-stone-400 hover:text-white hover:bg-white/10 rounded-lg transition-colors cursor-pointer"
                title="Zoom arrière"
                aria-label="Zoom arrière"
            >
                <Minus size={14} />
            </button>
            <button
                onclick={resetZoom}
                class="px-2.5 py-1 text-white hover:bg-white/10 rounded-lg text-xs font-bold font-mono transition-colors min-w-[55px] text-center cursor-pointer"
                title="Réinitialiser la vue"
                aria-label="Réinitialiser la vue"
            >
                {Math.round(zoom * 100)}%
            </button>
            <button
                onclick={zoomIn}
                class="p-1.5 text-stone-400 hover:text-white hover:bg-white/10 rounded-lg transition-colors cursor-pointer"
                title="Zoom avant"
                aria-label="Zoom avant"
            >
                <Plus size={14} />
            </button>
        </div>
    {:else}
        <!-- Premium Empty State -->
        <div class="absolute inset-0 flex flex-col items-center justify-center p-6 bg-stone-900 z-10">
            <div class="max-w-md w-full text-center space-y-6 bg-stone-900/60 backdrop-blur-md border border-stone-800/80 p-8 rounded-3xl shadow-2xl animate-in fade-in zoom-in duration-300">
                <div class="relative w-20 h-20 mx-auto flex items-center justify-center bg-burnt-orange/10 rounded-full border border-burnt-orange/20 text-burnt-orange animate-pulse">
                    <Compass size={40} />
                </div>
                <div class="space-y-2">
                    <h3 class="text-xl font-bold text-white tracking-wide">
                        {#if isGM}
                            Aucun plateau actif
                        {:else}
                            En attente du MJ
                        {/if}
                    </h3>
                    <p class="text-sm text-stone-400 leading-relaxed">
                        {#if isGM}
                            Créez un plateau et ajoutez-y une carte dans le panneau latéral droit (**onglet Plateaux**), puis activez-les pour lancer l'aventure.
                        {:else}
                            Le Maître du Jeu prépare le terrain. La carte s'affichera automatiquement ici dès qu'elle sera activée.
                        {/if}
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>
