<script lang="ts">
    import { activeBoardStore, sendPing, onPingReceived, onTokenDragged, sendTokenDrag } from "$lib/websocket";
    import { Compass, Minus, Plus, Trash2, Swords, Eye, EyeOff, Activity } from "lucide-svelte";
    import { fade } from "svelte/transition";
    import { onMount, untrack } from "svelte";
    import { page } from "$app/state";
    import { fetchPlayers, fetchGame } from "$lib/api";
    import {
        fetchBoardTokens,
        updateBoardTokenPosition,
        deleteBoardToken,
        updateTokenCombatState,
        type BoardToken
    } from "$lib/api/board";
    import { authClient } from "$lib/auth-client";
    import { parseConditions } from "$lib/api/character";
    import { supabase } from "$lib/supabaseClient";


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
    let currentUserId = $state("");

    // Tokens State
    let tokens = $state<BoardToken[]>([]);
    let draggedTokenId = $state<string | null>(null);
    let dragOffset = { x: 0, y: 0 };
    let tokenImages = $state<Record<string, HTMLImageElement>>({});
    let activeBoardTokensLoading = $state(true);
    const activeBoardId = $derived($activeBoardStore?.id || "");
    const pixelsPerCell = $derived($activeBoardStore?.pixels_per_cell || 70);
    const gridOffsetX = $derived($activeBoardStore?.grid_offset_x || 0);
    const gridOffsetY = $derived($activeBoardStore?.grid_offset_y || 0);
    // Token radius is 45% of a cell (leaves a 10% margin on each side)
    const tokenRadius = $derived(pixelsPerCell * 0.45);

    const combatActive = $derived($activeBoardStore?.combat_active || false);
    const activeTokenId = $derived($activeBoardStore?.combat_active_token_id || null);

    const visibleTokens = $derived(
        isGM ? tokens : tokens.filter(t => !t.is_hidden)
    );

    const AVAILABLE_CONDITIONS = [
        { id: "poisoned", dotColor: "#10B981" },
        { id: "stunned", dotColor: "#A855F7" },
        { id: "prone", dotColor: "#F59E0B" },
        { id: "blessed", dotColor: "#FACC15" },
        { id: "invisible", dotColor: "#6366F1" },
        { id: "unconscious", dotColor: "#BE123C" }
    ];

    // Context Menu State
    let contextMenu = $state<{
        x: number;
        y: number;
        tokenId: string;
        characterName: string;
        inCombat: boolean;
        isHidden: boolean;
        auraRadius?: number | null;
    } | null>(null);

    async function removeTokenFromBoard(tokenId: string) {
        // Optimistic update: remove locally immediately
        tokens = tokens.filter(t => t.id !== tokenId);
        draw();

        try {
            await deleteBoardToken(tokenId);
        } catch (error) {
            console.error("Failed to delete token:", error);
            // Restore token list if DB delete fails
            loadTokens();
        } finally {
            contextMenu = null;
        }
    }

    async function toggleTokenCombat(tokenId: string, currentInCombat: boolean) {
        try {
            const token = tokens.find(t => t.id === tokenId);
            let initVal = 0;
            if (!currentInCombat && token?.character?.is_npc) {
                const mod = token.character.initiative || 0;
                initVal = Math.floor(Math.random() * 20) + 1 + mod;
            }
            await updateTokenCombatState(tokenId, { in_combat: !currentInCombat, initiative: initVal });
        } catch (e) {
            console.error("Failed to toggle combat state from scene menu:", e);
        } finally {
            contextMenu = null;
        }
    }

    async function changeTokenFaction(tokenId: string, faction: 'ally' | 'enemy' | 'neutral') {
        try {
            // Optimistic update
            tokens = tokens.map(t => t.id === tokenId ? { ...t, faction } : t);
            draw();
            await updateTokenCombatState(tokenId, { faction });
        } catch (e) {
            console.error("Failed to update token faction from scene menu:", e);
            loadTokens();
        } finally {
            contextMenu = null;
        }
    }

    async function toggleTokenVisibility(tokenId: string, currentIsHidden: boolean) {
        try {
            await updateTokenCombatState(tokenId, { is_hidden: !currentIsHidden });
        } catch (e) {
            console.error("Failed to toggle visibility state from scene menu:", e);
        } finally {
            contextMenu = null;
        }
    }

    async function promptCircleDiameter(tokenId: string, currentRadius: number | null) {
        contextMenu = null;
        const currentDiameter = currentRadius ? currentRadius * 2 : 0;
        const msg = currentDiameter > 0 
            ? `Entrez le diamètre du cercle en mètres (actuel : ${currentDiameter}m, ou 0 pour le supprimer) :`
            : "Entrez le diamètre du cercle à dessiner autour du jeton (en mètres, ex : 6) :";
        
        const input = prompt(msg, currentDiameter > 0 ? String(currentDiameter) : "6");
        if (input === null) return;

        const val = parseFloat(input);
        if (isNaN(val) || val < 0) {
            alert("Veuillez entrer un nombre positif valide.");
            return;
        }

        const newRadius = val > 0 ? val / 2 : null;
        try {
            // Optimistic update
            tokens = tokens.map(t => t.id === tokenId ? { ...t, aura_radius: newRadius } : t);
            draw();
            await updateTokenCombatState(tokenId, { aura_radius: newRadius });
        } catch (e) {
            console.error("Failed to update token aura radius:", e);
            loadTokens();
        }
    }

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

    // Preload token images when tokens change
    $effect(() => {
        tokens.forEach(token => {
            const url = token.character?.avatar_url;
            if (url && !tokenImages[url]) {
                const imgObj = new Image();
                imgObj.src = url;
                imgObj.onload = () => {
                    tokenImages = { ...tokenImages, [url]: imgObj };
                    draw();
                };
            }
        });
    });

    // Continuous rendering during combat for smooth pulsing active turn ring
    $effect(() => {
        if (combatActive) {
            let activeAnimId: number;
            const tick = () => {
                draw();
                activeAnimId = requestAnimationFrame(tick);
            };
            activeAnimId = requestAnimationFrame(tick);
            return () => cancelAnimationFrame(activeAnimId);
        }
    });

    async function loadTokens() {
        if (!activeBoardId) {
            tokens = [];
            return;
        }
        try {
            activeBoardTokensLoading = true;
            tokens = await fetchBoardTokens(activeBoardId);
            draw();
        } catch (error) {
            console.error("Failed to fetch board tokens:", error);
        } finally {
            activeBoardTokensLoading = false;
        }
    }

    let tokensChannel: any = null;
    let charactersChannel: any = null;
    $effect(() => {
        if (activeBoardId) {
            loadTokens();

            if (tokensChannel) {
                supabase.removeChannel(tokensChannel);
            }
            if (charactersChannel) {
                supabase.removeChannel(charactersChannel);
            }

            tokensChannel = supabase.channel(`board_tokens:${activeBoardId}`)
                .on(
                    'postgres_changes',
                    { event: '*', schema: 'public', table: 'board_tokens' },
                    async (payload) => {
                        if (payload.eventType === 'INSERT') {
                            if (payload.new.board_id === activeBoardId) {
                                loadTokens();
                            }
                        } else if (payload.eventType === 'DELETE') {
                            // Check if deleted token was in our current tokens
                            if (tokens.some(t => t.id === payload.old.id)) {
                                tokens = tokens.filter(t => t.id !== payload.old.id);
                                draw();
                            }
                        } else if (payload.eventType === 'UPDATE') {
                            if (payload.new.board_id === activeBoardId) {
                                const updatedToken = payload.new;
                                tokens = tokens.map(t => {
                                    if (t.id === updatedToken.id) {
                                        if (draggedTokenId === updatedToken.id) {
                                            return { ...t, ...updatedToken, x: t.x, y: t.y };
                                        } else {
                                            return { ...t, ...updatedToken };
                                        }
                                    }
                                    return t;
                                });
                                draw();
                            }
                        }
                    }
                )
                .subscribe();

            const gameId = page.params.id || "";
            charactersChannel = supabase.channel(`board_characters:${gameId}`)
                .on(
                    'postgres_changes',
                    { event: 'UPDATE', schema: 'public', table: 'characters' },
                    (payload) => {
                        const updatedChar = payload.new;
                        tokens = tokens.map(t => {
                            if (t.character_id === updatedChar.id) {
                                return {
                                    ...t,
                                    character: t.character ? {
                                        ...t.character,
                                        ...updatedChar
                                    } : undefined
                                };
                            }
                            return t;
                        });
                        draw();
                    }
                )
                .subscribe();
        } else {
            tokens = [];
        }

        return () => {
            if (tokensChannel) {
                supabase.removeChannel(tokensChannel);
                tokensChannel = null;
            }
            if (charactersChannel) {
                supabase.removeChannel(charactersChannel);
                charactersChannel = null;
            }
        };
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
                    currentUserId = sessionData.user.id;
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

        const unsubscribeTokenDrag = onTokenDragged((dragData) => {
            if (draggedTokenId !== dragData.tokenId) {
                tokens = tokens.map(t => {
                    if (t.id === dragData.tokenId) {
                        return { ...t, x: dragData.x, y: dragData.y };
                    }
                    return t;
                });
                draw();
            }
        });

        const handleWindowMouseDown = (e: MouseEvent) => {
            if (contextMenu && containerEl && !containerEl.contains(e.target as Node)) {
                contextMenu = null;
            }
        };

        window.addEventListener("mousedown", handleWindowMouseDown);

        return () => {
            unsubscribePing();
            unsubscribeTokenDrag();
            window.removeEventListener("mousedown", handleWindowMouseDown);
            if (containerEl) {
                containerEl.removeEventListener("wheel", handleWheel);
            }
            if (animId !== null) {
                cancelAnimationFrame(animId);
            }
        };
    });

    function handleMouseDown(e: MouseEvent) {
        if (!activeMapUrl || !imgLoaded || !width || !height) return;

        if (e.button === 0 || e.button === 1) {
            contextMenu = null;
        }

        if (e.button === 1) {
            // Middle-click to pan
            e.preventDefault();
            isPanning = true;
            startX = e.clientX;
            startY = e.clientY;
            startPanX = pan.x;
            startPanY = pan.y;
        } else if (e.button === 0) {
            // Left-click: check if clicked on a token first
            if (containerEl) {
                const rect = containerEl.getBoundingClientRect();
                const mouseX = e.clientX - rect.left;
                const mouseY = e.clientY - rect.top;

                // Convert screen mouse coordinate to image space coordinate
                const x_image_centered = (mouseX - pan.x) / zoom;
                const y_image_centered = (mouseY - pan.y) / zoom;

                const halfW = naturalWidth / 2;
                const halfH = naturalHeight / 2;
                const tx = x_image_centered + halfW;
                const ty = y_image_centered + halfH;

                // Find clicked token (reverse loop to get topmost first)
                let clickedToken: BoardToken | null = null;
                for (let i = visibleTokens.length - 1; i >= 0; i--) {
                    const t = visibleTokens[i];
                    const tokenX = t.x * naturalWidth;
                    const tokenY = t.y * naturalHeight;
                    const dist = Math.sqrt((tx - tokenX) ** 2 + (ty - tokenY) ** 2);
                    if (dist <= tokenRadius) {
                        clickedToken = t;
                        break;
                    }
                }

                if (clickedToken) {
                    // Check permissions: GM can move all, Player can move own character
                    const canMove = isGM || (clickedToken.character?.user_id === currentUserId && !clickedToken.character?.is_npc);
                    if (canMove) {
                        draggedTokenId = clickedToken.id;
                        dragOffset.x = tx - (clickedToken.x * naturalWidth);
                        dragOffset.y = ty - (clickedToken.y * naturalHeight);
                    } else {
                        console.warn("No permission to move this token");
                    }
                    return; // Stop here, do not ping
                }
            }

            // Normal left-click to ping
            handleSceneClick(e);
        }
    }

    function handleMouseMove(e: MouseEvent) {
        if (isPanning) {
            const dx = e.clientX - startX;
            const dy = e.clientY - startY;
            pan.x = startPanX + dx;
            pan.y = startPanY + dy;
            
            if (animId === null) {
                draw();
            }
        } else if (draggedTokenId && containerEl) {
            const rect = containerEl.getBoundingClientRect();
            const mouseX = e.clientX - rect.left;
            const mouseY = e.clientY - rect.top;

            const x_image_centered = (mouseX - pan.x) / zoom;
            const y_image_centered = (mouseY - pan.y) / zoom;

            const halfW = naturalWidth / 2;
            const halfH = naturalHeight / 2;
            const tx = x_image_centered + halfW;
            const ty = y_image_centered + halfH;

            // Clamped coordinates relative to image bounds
            const finalTx = Math.min(Math.max(tx - dragOffset.x, 0), naturalWidth);
            const finalTy = Math.min(Math.max(ty - dragOffset.y, 0), naturalHeight);

            const normX = finalTx / naturalWidth;
            const normY = finalTy / naturalHeight;

            // Update locally for smooth rendering
            tokens = tokens.map(t => {
                if (t.id === draggedTokenId) {
                    return { ...t, x: normX, y: normY };
                }
                return t;
            });
            draw();

            // Broadcast real-time drag position
            sendTokenDrag(draggedTokenId, normX, normY);
        }
    }

    async function handleMouseUp(e: MouseEvent) {
        if (e.button === 1) {
            isPanning = false;
        } else if (e.button === 0 && draggedTokenId) {
            // Drag end: save coordinates to DB
            const draggedToken = tokens.find(t => t.id === draggedTokenId);
            if (draggedToken) {
                try {
                    await updateBoardTokenPosition(draggedTokenId, draggedToken.x, draggedToken.y);
                } catch (error) {
                    console.error("Failed to save token position to DB:", error);
                }
            }
            draggedTokenId = null;
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

        const halfW = naturalWidth / 2;
        const halfH = naturalHeight / 2;

        // Draw subtle border around map
        ctx.strokeStyle = "rgba(0, 0, 0, 0.25)";
        ctx.lineWidth = 1 / zoom;
        ctx.strokeRect(-halfW, -halfH, naturalWidth, naturalHeight);

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

        // Draw active tokens
        for (const token of visibleTokens) {
            const tx = token.x * naturalWidth - halfW;
            const ty = token.y * naturalHeight - halfH;
            const radius = tokenRadius;

            const isMonster = token.character?.type === "MONSTER";
            const isCharacterInvisible = token.character?.conditions?.includes("invisible") || false;

            // RPG rule: If character is invisible, and it is a monster and we are not the GM, we skip drawing it completely.
            if (isCharacterInvisible && isMonster && !isGM) {
                continue;
            }

            // Draw Aura/Circle Underneath
            if (token.aura_radius && token.aura_radius > 0) {
                const auraRadiusPx = token.aura_radius * pixelsPerCell;
                ctx.save();
                
                let auraColor = "#3B82F6"; // default blue for ally
                if (token.faction === "enemy") {
                    auraColor = "#EF4444"; // red for enemy
                } else if (token.faction === "neutral") {
                    auraColor = "#EAB308"; // yellow for neutral
                }
                
                // Soft filled background
                ctx.globalAlpha = 0.08;
                ctx.fillStyle = auraColor;
                ctx.beginPath();
                ctx.arc(tx, ty, auraRadiusPx, 0, Math.PI * 2);
                ctx.fill();
                
                // Dashed outline
                ctx.globalAlpha = 0.5;
                ctx.strokeStyle = auraColor;
                ctx.lineWidth = 2 / zoom;
                ctx.setLineDash([6 / zoom, 4 / zoom]);
                ctx.beginPath();
                ctx.arc(tx, ty, auraRadiusPx, 0, Math.PI * 2);
                ctx.stroke();
                
                ctx.restore();
            }

            ctx.save();
            
            // Set transparency for invisible tokens
            if (isCharacterInvisible) {
                ctx.globalAlpha = 0.5;
            }

            // 0. Draw Active Combat Turn ring underneath
            if (combatActive && token.id === activeTokenId) {
                ctx.save();
                ctx.beginPath();
                ctx.arc(tx, ty, radius + 5 / zoom, 0, Math.PI * 2);
                ctx.strokeStyle = "#E07A5F"; // theme burnt orange
                ctx.lineWidth = 5 / zoom;
                
                // Pulsing glow effect
                const pulse = Math.abs(Math.sin(Date.now() / 200)) * 3 / zoom;
                ctx.shadowBlur = 10 + pulse * zoom;
                ctx.shadowColor = "#E07A5F";
                ctx.stroke();
                ctx.restore();
            }

            // 1. Draw Ring / Border
            ctx.save();
            ctx.beginPath();
            ctx.arc(tx, ty, radius + 2 / zoom, 0, Math.PI * 2);

            let color = "#3B82F6"; // default blue for ally
            if (token.faction === "enemy") {
                color = "#EF4444"; // red for enemy
            } else if (token.faction === "neutral") {
                color = "#EAB308"; // yellow for neutral
            }

            if (token.character?.user_id === currentUserId && !token.character?.is_npc) {
                color = "#E07A5F"; // orange for self (only for player characters)
            }

            // If token is hidden from players and we are GM, color it grey or dash it
            if (token.is_hidden && isGM) {
                color = "#78716c"; // stone-500
                ctx.setLineDash([5 / zoom, 5 / zoom]);
            }

            ctx.strokeStyle = color;
            ctx.lineWidth = 3 / zoom;
            ctx.stroke();
            ctx.restore();

            // 2. Draw Avatar
            const avatarUrl = token.character?.avatar_url;
            const loadedImg = avatarUrl ? tokenImages[avatarUrl] : null;

            ctx.save();
            ctx.beginPath();
            ctx.arc(tx, ty, radius, 0, Math.PI * 2);
            ctx.clip();

            if (loadedImg) {
                ctx.drawImage(loadedImg, tx - radius, ty - radius, radius * 2, radius * 2);
            } else {
                const grad = ctx.createRadialGradient(tx, ty, 2, tx, ty, radius);
                grad.addColorStop(0, "#f5f5f4");
                grad.addColorStop(1, "#d6d3d1");
                ctx.fillStyle = grad;
                ctx.fill();

                ctx.fillStyle = "#57534e";
                ctx.font = `bold ${Math.max(10, 12 / zoom)}px sans-serif`;
                ctx.textAlign = "center";
                ctx.textBaseline = "middle";
                ctx.fillText((token.character?.name || "?").charAt(0).toUpperCase(), tx, ty);
            }
            ctx.restore();

            // 3. Draw Condition Dots (top-right)
            const conditions = parseConditions(token.character?.conditions);
            if (conditions.length > 0) {
                ctx.save();
                const dotRadius = 4 / zoom;
                const gap = 10 / zoom;
                let dotX = tx + radius - 6 / zoom;
                let dotY = ty - radius + 6 / zoom;

                conditions.forEach((condId, idx) => {
                    const predefined = AVAILABLE_CONDITIONS.find(c => c.id === condId);
                    ctx.beginPath();
                    ctx.arc(dotX - idx * gap, dotY, dotRadius, 0, Math.PI * 2);
                    ctx.fillStyle = predefined ? predefined.dotColor : "#A8A29E"; // stone-400 for custom
                    ctx.fill();
                    ctx.strokeStyle = "#ffffff";
                    ctx.lineWidth = 1 / zoom;
                    ctx.stroke();
                });
                ctx.restore();
            }

            // 4. Draw Name Label Under Token
            ctx.save();
            let labelText = token.character?.name || "Sans nom";
            if (token.is_hidden && isGM) {
                labelText = `👁️ ${labelText}`;
            }
            ctx.font = `bold ${10 / zoom}px sans-serif`;
            const textWidth = ctx.measureText(labelText).width;
            const labelW = textWidth + 8 / zoom;
            const labelH = 14 / zoom;
            const labelX = tx - labelW / 2;
            const labelY = ty + radius + 4 / zoom;

            ctx.fillStyle = "rgba(0, 0, 0, 0.7)";
            ctx.beginPath();
            const r = 3 / zoom;
            if (ctx.roundRect) {
                ctx.roundRect(labelX, labelY, labelW, labelH, r);
            } else {
                ctx.rect(labelX, labelY, labelW, labelH);
            }
            ctx.fill();

            ctx.fillStyle = "#ffffff";
            ctx.textAlign = "center";
            ctx.textBaseline = "middle";
            ctx.fillText(labelText, tx, labelY + labelH / 2);
            ctx.restore();

            ctx.restore(); // Restore globalAlpha
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

    function handleContextMenu(e: MouseEvent) {
        if (!activeMapUrl || !imgLoaded || !width || !height) return;

        if (containerEl) {
            const rect = containerEl.getBoundingClientRect();
            const mouseX = e.clientX - rect.left;
            const mouseY = e.clientY - rect.top;

            const x_image_centered = (mouseX - pan.x) / zoom;
            const y_image_centered = (mouseY - pan.y) / zoom;

            const halfW = naturalWidth / 2;
            const halfH = naturalHeight / 2;
            const tx = x_image_centered + halfW;
            const ty = y_image_centered + halfH;

            // Find clicked token (reverse loop to get topmost first)
            let clickedToken: BoardToken | null = null;
            for (let i = visibleTokens.length - 1; i >= 0; i--) {
                const t = visibleTokens[i];
                const tokenX = t.x * naturalWidth;
                const tokenY = t.y * naturalHeight;
                const dist = Math.sqrt((tx - tokenX) ** 2 + (ty - tokenY) ** 2);
                if (dist <= tokenRadius) {
                    clickedToken = t;
                    break;
                }
            }

            if (clickedToken && isGM) {
                // Right-click on a token: show custom menu, prevent browser menu
                e.preventDefault();
                contextMenu = {
                    x: mouseX,
                    y: mouseY,
                    tokenId: clickedToken.id,
                    characterName: clickedToken.character?.name || "Personnage",
                    inCombat: clickedToken.in_combat,
                    isHidden: clickedToken.is_hidden,
                    auraRadius: clickedToken.aura_radius
                };
            } else {
                // Right-click on empty map: do not prevent default, just close custom menu
                contextMenu = null;
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
    oncontextmenu={handleContextMenu}
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

        <!-- Custom Context Menu for tokens -->
        {#if contextMenu && isGM}
            <div
                role="menu"
                tabindex="-1"
                class="absolute bg-white rounded-xl shadow-xl border border-stone-200 py-1 z-40 w-44 animate-in fade-in zoom-in-95 duration-100 origin-top-left flex flex-col"
                style="left: {contextMenu.x}px; top: {contextMenu.y}px;"
                onmousedown={(e) => e.stopPropagation()}
            >
                <div class="px-3 py-1 text-[10px] font-bold text-stone-400 uppercase tracking-wider border-b border-stone-100 mb-1 truncate" title={contextMenu.characterName}>
                    {contextMenu.characterName}
                </div>
                <button
                    role="menuitem"
                    onclick={() => toggleTokenCombat(contextMenu?.tokenId || '', contextMenu?.inCombat || false)}
                    class="w-full px-3 py-2 text-left text-xs text-stone-700 hover:bg-stone-50 font-semibold flex items-center gap-2 transition-colors cursor-pointer"
                >
                    <Swords size={13} class="text-stone-400" />
                    {contextMenu.inCombat ? 'Retirer du combat' : 'Ajouter au combat'}
                </button>
                <button
                    role="menuitem"
                    onclick={() => toggleTokenVisibility(contextMenu?.tokenId || '', contextMenu?.isHidden || false)}
                    class="w-full px-3 py-2 text-left text-xs text-stone-700 hover:bg-stone-50 font-semibold flex items-center gap-2 transition-colors cursor-pointer"
                >
                    {#if contextMenu.isHidden}
                        <Eye size={13} class="text-stone-400" />
                        Rendre visible
                    {:else}
                        <EyeOff size={13} class="text-stone-400" />
                        Masquer aux joueurs
                    {/if}
                </button>
                <button
                    role="menuitem"
                    onclick={() => promptCircleDiameter(contextMenu?.tokenId || '', contextMenu?.auraRadius || null)}
                    class="w-full px-3 py-2 text-left text-xs text-stone-700 hover:bg-stone-50 font-semibold flex items-center gap-2 transition-colors cursor-pointer"
                >
                    <Activity size={13} class="text-stone-400" />
                    <span>Dessiner un cercle</span>
                </button>
                <div class="h-px bg-stone-100 my-1"></div>
                <div class="px-3 py-1 text-[9px] font-bold text-stone-400 uppercase tracking-wider">Faction</div>
                <div class="flex px-2 py-1 gap-1">
                    <button
                        onclick={() => changeTokenFaction(contextMenu?.tokenId || '', 'ally')}
                        class="flex-1 py-1 rounded text-[10px] font-bold text-center transition-colors bg-emerald-50 text-emerald-600 hover:bg-emerald-100 cursor-pointer"
                    >
                        Allié
                    </button>
                    <button
                        onclick={() => changeTokenFaction(contextMenu?.tokenId || '', 'enemy')}
                        class="flex-1 py-1 rounded text-[10px] font-bold text-center transition-colors bg-rose-50 text-rose-600 hover:bg-rose-100 cursor-pointer"
                    >
                        Ennemi
                    </button>
                    <button
                        onclick={() => changeTokenFaction(contextMenu?.tokenId || '', 'neutral')}
                        class="flex-1 py-1 rounded text-[10px] font-bold text-center transition-colors bg-amber-50 text-amber-600 hover:bg-amber-100 cursor-pointer"
                    >
                        Neutre
                    </button>
                </div>
                <div class="h-px bg-stone-100 my-1"></div>
                <button
                    role="menuitem"
                    onclick={() => removeTokenFromBoard(contextMenu?.tokenId || '')}
                    class="w-full px-3 py-2 text-left text-xs text-red-500 hover:bg-red-50 font-bold flex items-center gap-2 transition-colors cursor-pointer"
                >
                    <Trash2 size={13} />
                    Retirer du plateau
                </button>
            </div>
        {/if}
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
                            Créez un plateau et ajoutez-y une carte dans le panneau latéral droit (<strong>onglet Plateaux</strong>), puis activez-les pour lancer l'aventure.
                        {:else}
                            Le Maître du Jeu prépare le terrain. La carte s'affichera automatiquement ici dès qu'elle sera activée.
                        {/if}
                    </p>
                </div>
            </div>
        </div>
    {/if}
</div>
