<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import Header from "$lib/components/Header.svelte";
    import { onMount } from "svelte";
    import { fetchBoard, updateBoard } from "$lib/api/board";
    import { uploadImage } from "$lib/api/storage";
    import { authClient } from "$lib/auth-client";
    import {
        ArrowLeft, Check, Upload, Link as LinkIcon,
        Loader2, Image as ImageIcon, Layers, Grid3x3,
        ZoomIn, ZoomOut
    } from "lucide-svelte";

    const gameId = page.params.id as string;
    const boardId = page.params.boardId as string;

    let boardName = $state("");
    let boardUrl = $state("");
    let uploadMethod = $state<"file" | "url">("file");

    // Scale fields
    let pixelsPerCell = $state(70);
    let gridOffsetX = $state(0);
    let gridOffsetY = $state(0);

    let loading = $state(true);
    let saving = $state(false);
    let isUploadingFile = $state(false);
    let error = $state<string | null>(null);

    let fileInput = $state<HTMLInputElement | null>(null);
    let dragOver = $state(false);

    // Calibration canvas state
    let calibrationCanvas = $state<HTMLCanvasElement | null>(null);
    let calibrationContainer = $state<HTMLElement | null>(null);
    let calibrationImg: HTMLImageElement | undefined = undefined;
    let calibrationImgLoaded = $state(false);
    let calZoom = $state(1);
    let calPan = { x: 0, y: 0 };
    let calIsPanning = $state(false);
    let calStartX = 0;
    let calStartY = 0;
    let calStartPanX = 0;
    let calStartPanY = 0;
    let calWidth = $state(0);
    let calHeight = $state(0);

    onMount(async () => {
        try {
            const { data: tokenData } = await authClient.token();
            if (tokenData?.token) {
                const boardData = await fetchBoard(boardId);
                if (boardData) {
                    boardName = boardData.name;
                    boardUrl = boardData.image_url || "";
                    pixelsPerCell = boardData.pixels_per_cell ?? 70;
                    gridOffsetX = boardData.grid_offset_x ?? 0;
                    gridOffsetY = boardData.grid_offset_y ?? 0;
                    if (boardUrl && !boardUrl.startsWith("http")) {
                        uploadMethod = "file";
                    } else if (boardUrl) {
                        // Default to URL if it looks like an absolute external URL, or fallback to file
                        uploadMethod = boardUrl.startsWith("http") && !boardUrl.includes("supabase") ? "url" : "file";
                    }
                } else {
                    error = "Plateau introuvable.";
                }
            }
        } catch (e) {
            console.error("Failed to fetch board details:", e);
            error = "Impossible de charger les informations du plateau.";
        } finally {
            loading = false;
        }
    });

    // Load calibration image when boardUrl changes
    $effect(() => {
        if (boardUrl) {
            calibrationImgLoaded = false;
            calibrationImg = new Image();
            calibrationImg.crossOrigin = "anonymous";
            calibrationImg.src = boardUrl;
            calibrationImg.onload = () => {
                calibrationImgLoaded = true;
                if (calibrationImg && calWidth && calHeight) {
                    const fitScale = Math.min(
                        (calWidth * 0.95) / calibrationImg.naturalWidth,
                        (calHeight * 0.95) / calibrationImg.naturalHeight
                    );
                    calZoom = Math.min(Math.max(fitScale, 0.05), 2);
                    calPan.x = calWidth / 2;
                    calPan.y = calHeight / 2;
                }
                drawCalibration();
            };
        } else {
            calibrationImgLoaded = false;
        }
    });

    // Redraw calibration when params change
    $effect(() => {
        // Track reactive deps
        const _ = [pixelsPerCell, gridOffsetX, gridOffsetY, calZoom, calWidth, calHeight];
        if (calibrationImgLoaded) {
            drawCalibration();
        }
    });

    // Set up wheel handler on mount
    $effect(() => {
        if (calibrationContainer) {
            const handler = (e: WheelEvent) => {
                e.preventDefault();
                if (e.ctrlKey) {
                    const factor = e.deltaY < 0 ? 1.08 : 1 / 1.08;
                    const newZoom = Math.min(Math.max(calZoom * factor, 0.05), 8);
                    if (newZoom !== calZoom) {
                        const rect = calibrationContainer!.getBoundingClientRect();
                        const mouseX = e.clientX - rect.left;
                        const mouseY = e.clientY - rect.top;
                        calPan.x = mouseX - (mouseX - calPan.x) * (newZoom / calZoom);
                        calPan.y = mouseY - (mouseY - calPan.y) * (newZoom / calZoom);
                        calZoom = newZoom;
                    }
                } else {
                    if (e.shiftKey) {
                        calPan.x -= e.deltaY;
                    } else {
                        calPan.x -= e.deltaX;
                        calPan.y -= e.deltaY;
                    }
                }
                drawCalibration();
            };
            calibrationContainer.addEventListener("wheel", handler, { passive: false });
            return () => calibrationContainer?.removeEventListener("wheel", handler);
        }
    });

    function drawCalibration() {
        if (!calibrationCanvas || !calibrationImgLoaded || !calibrationImg || !calWidth || !calHeight) return;
        const ctx = calibrationCanvas.getContext("2d");
        if (!ctx) return;

        const dpr = window.devicePixelRatio || 1;
        ctx.clearRect(0, 0, calWidth * dpr, calHeight * dpr);

        ctx.save();
        ctx.scale(dpr, dpr);
        ctx.translate(calPan.x, calPan.y);
        ctx.scale(calZoom, calZoom);

        const nw = calibrationImg.naturalWidth;
        const nh = calibrationImg.naturalHeight;

        // Draw image centered
        ctx.drawImage(calibrationImg, -nw / 2, -nh / 2, nw, nh);

        const halfW = nw / 2;
        const halfH = nh / 2;

        // Draw grid lines
        if (pixelsPerCell > 0) {
            ctx.strokeStyle = "rgba(59, 130, 246, 0.45)";
            ctx.lineWidth = 1 / calZoom;

            const oX = gridOffsetX - halfW;
            const oY = gridOffsetY - halfH;

            // Vertical lines
            const startColPx = Math.floor((-halfW - oX) / pixelsPerCell) * pixelsPerCell + oX;
            for (let x = startColPx; x <= halfW; x += pixelsPerCell) {
                ctx.beginPath();
                ctx.moveTo(x, -halfH);
                ctx.lineTo(x, halfH);
                ctx.stroke();
            }

            // Horizontal lines
            const startRowPx = Math.floor((-halfH - oY) / pixelsPerCell) * pixelsPerCell + oY;
            for (let y = startRowPx; y <= halfH; y += pixelsPerCell) {
                ctx.beginPath();
                ctx.moveTo(-halfW, y);
                ctx.lineTo(halfW, y);
                ctx.stroke();
            }
        }

        // Border
        ctx.strokeStyle = "rgba(0, 0, 0, 0.3)";
        ctx.lineWidth = 2 / calZoom;
        ctx.strokeRect(-halfW, -halfH, nw, nh);

        ctx.restore();
    }

    function handleCalMouseDown(e: MouseEvent) {
        if (e.button === 0 || e.button === 1) {
            e.preventDefault();
            calIsPanning = true;
            calStartX = e.clientX;
            calStartY = e.clientY;
            calStartPanX = calPan.x;
            calStartPanY = calPan.y;
        }
    }

    function handleCalMouseMove(e: MouseEvent) {
        if (calIsPanning) {
            calPan.x = calStartPanX + (e.clientX - calStartX);
            calPan.y = calStartPanY + (e.clientY - calStartY);
            drawCalibration();
        }
    }

    function handleCalMouseUp() {
        calIsPanning = false;
    }

    function calZoomIn() {
        const newZoom = Math.min(calZoom * 1.3, 8);
        if (calWidth && calHeight) {
            const cx = calWidth / 2;
            const cy = calHeight / 2;
            calPan.x = cx - (cx - calPan.x) * (newZoom / calZoom);
            calPan.y = cy - (cy - calPan.y) * (newZoom / calZoom);
        }
        calZoom = newZoom;
        drawCalibration();
    }

    function calZoomOut() {
        const newZoom = Math.max(calZoom / 1.3, 0.05);
        if (calWidth && calHeight) {
            const cx = calWidth / 2;
            const cy = calHeight / 2;
            calPan.x = cx - (cx - calPan.x) * (newZoom / calZoom);
            calPan.y = cy - (cy - calPan.y) * (newZoom / calZoom);
        }
        calZoom = newZoom;
        drawCalibration();
    }

    function calResetView() {
        if (calibrationImg && calWidth && calHeight) {
            const fitScale = Math.min(
                (calWidth * 0.95) / calibrationImg.naturalWidth,
                (calHeight * 0.95) / calibrationImg.naturalHeight
            );
            calZoom = Math.min(Math.max(fitScale, 0.05), 2);
            calPan.x = calWidth / 2;
            calPan.y = calHeight / 2;
            drawCalibration();
        }
    }

    async function handleFileUpload(file: File) {
        if (isUploadingFile) return;
        try {
            isUploadingFile = true;
            const publicUrl = await uploadImage(file);
            boardUrl = publicUrl;
        } catch (err) {
            console.error("Erreur lors de l'upload de l'image :", err);
            alert("Erreur lors de l'upload de l'image.");
        } finally {
            isUploadingFile = false;
        }
    }

    function onFileSelected(e: Event) {
        const files = (e.target as HTMLInputElement).files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
    }

    function handleDragOver(e: DragEvent) {
        e.preventDefault();
        dragOver = true;
    }

    function handleDragLeave() {
        dragOver = false;
    }

    function handleDrop(e: DragEvent) {
        e.preventDefault();
        dragOver = false;
        const files = e.dataTransfer?.files;
        if (files && files.length > 0) {
            handleFileUpload(files[0]);
        }
    }

    async function handleSave() {
        if (!boardName.trim() || !boardUrl.trim()) return;
        saving = true;
        try {
            await updateBoard(
                boardId,
                boardName.trim(),
                boardUrl.trim(),
                pixelsPerCell,
                gridOffsetX,
                gridOffsetY
            );
            goto(`/table/${gameId}/gm/settings?tab=boards`);
        } catch (e) {
            console.error("Failed to save board settings:", e);
            alert("Erreur lors de la sauvegarde du plateau.");
        } finally {
            saving = false;
        }
    }
</script>

<svelte:window onmousemove={handleCalMouseMove} onmouseup={handleCalMouseUp} />

<div class="min-h-screen bg-cream">
    <Header />

    <main class="max-w-2xl mx-auto p-6 md:p-12 animate-in fade-in duration-300">
        <!-- Back Button -->
        <a
            href="/table/{gameId}/gm/settings?tab=boards"
            class="inline-flex items-center gap-2 text-stone-500 hover:text-dark-gray mb-6 font-medium transition-colors"
        >
            <ArrowLeft size={16} />
            Retour aux paramètres
        </a>

        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-display font-bold text-dark-gray mb-2">
                Modifier le plateau
            </h1>
            <p class="text-stone-500">
                Ajustez le nom, l'image et l'échelle de la carte pour ce plateau de jeu.
            </p>
        </div>

        {#if loading}
            <div class="flex justify-center py-12">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-burnt-orange"></div>
            </div>
        {:else if error}
            <div class="p-4 bg-red-50 text-red-600 rounded-xl border border-red-100 shadow-sm">
                {error}
            </div>
        {:else}
            <div class="bg-white rounded-2xl shadow-sm border border-stone-100 p-6 md:p-8 space-y-6">
                <!-- Header Icon & Title -->
                <div class="flex items-center gap-4 pb-6 border-b border-stone-100">
                    <div class="w-12 h-12 rounded-full bg-burnt-orange/10 flex items-center justify-center text-burnt-orange">
                        <Layers size={24} />
                    </div>
                    <div>
                        <h2 class="text-xl font-bold text-dark-gray">Paramètres de la carte</h2>
                        <p class="text-sm text-stone-400">Plateau ID: {boardId}</p>
                    </div>
                </div>

                <!-- Form Section -->
                <div class="space-y-4">
                    <div class="space-y-2">
                        <label for="board-name" class="text-sm font-bold text-dark-gray">Nom du plateau</label>
                        <input
                            id="board-name"
                            type="text"
                            placeholder="Ex: Donjon de l'effroi..."
                            bind:value={boardName}
                            required
                            class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all font-medium text-dark-gray"
                        />
                    </div>

                    <!-- Image Source Selection Swatches -->
                    <div class="space-y-2">
                        <span class="text-sm font-bold text-dark-gray block">Source de l'image</span>
                        <div class="flex rounded-lg bg-stone-100 p-0.5 gap-0.5 h-[42px] items-center max-w-sm">
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "file"; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full cursor-pointer
                                {uploadMethod === 'file' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <Upload size={14} />
                                Fichier local
                            </button>
                            <button
                                type="button"
                                onclick={() => { uploadMethod = "url"; }}
                                class="flex-1 py-1.5 text-xs font-bold rounded-md transition-all flex items-center justify-center gap-1.5 h-full cursor-pointer
                                {uploadMethod === 'url' ? 'bg-white text-stone-800 shadow-sm' : 'text-stone-500 hover:text-stone-700'}"
                            >
                                <LinkIcon size={14} />
                                Lien externe URL
                            </button>
                        </div>
                    </div>

                    <!-- Map Image Upload/Url Area -->
                    {#if uploadMethod === 'file'}
                        <div
                            role="button"
                            tabindex="0"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                            onclick={() => fileInput?.click()}
                            onkeydown={(e) => e.key === "Enter" && fileInput?.click()}
                            class="border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all flex flex-col items-center justify-center gap-3 bg-stone-50
                            {dragOver ? 'border-burnt-orange bg-burnt-orange/5' : 'border-stone-200 hover:border-stone-300'}"
                        >
                            <input
                                type="file"
                                accept="image/*"
                                bind:this={fileInput}
                                onchange={onFileSelected}
                                class="hidden"
                            />
                            {#if isUploadingFile}
                                <Loader2 size={28} class="animate-spin text-burnt-orange" />
                                <span class="text-sm text-stone-500 font-medium">Téléversement de la nouvelle carte...</span>
                            {:else if boardUrl}
                                <Check size={28} class="text-emerald-500" />
                                <span class="text-sm text-emerald-600 font-bold">Image chargée avec succès !</span>
                            {:else}
                                <Upload size={28} class="text-stone-400" />
                                <span class="text-sm text-stone-500">
                                    Glissez la nouvelle image ici ou <strong class="text-burnt-orange">parcourez vos fichiers</strong>
                                </span>
                            {/if}
                        </div>
                    {:else}
                        <div class="space-y-2">
                            <label for="board-url" class="text-sm font-bold text-dark-gray">Adresse URL de la carte</label>
                            <input
                                id="board-url"
                                type="url"
                                placeholder="https://exemple.com/carte.png"
                                bind:value={boardUrl}
                                required
                                class="w-full px-4 py-2.5 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-burnt-orange/20 focus:border-burnt-orange transition-all font-medium"
                            />
                        </div>
                    {/if}

                    <!-- Preview Section -->
                    {#if boardUrl}
                        <div class="space-y-2 pt-2">
                            <span class="text-sm font-bold text-dark-gray block">Aperçu de la carte</span>
                            <div class="relative rounded-2xl overflow-hidden border border-stone-200/80 shadow-md bg-stone-100 max-h-64 flex items-center justify-center">
                                <img 
                                    src={boardUrl} 
                                    alt="Aperçu du plateau"
                                    class="w-full h-full object-contain max-h-60"
                                />
                            </div>
                        </div>
                    {/if}
                </div>

                <!-- Scale Section -->
                {#if boardUrl}
                    <div class="space-y-5 pt-2 border-t border-stone-100">
                        <!-- Scale Header -->
                        <div class="flex items-center gap-3 pt-4">
                            <div class="w-10 h-10 rounded-full bg-blue-500/10 flex items-center justify-center text-blue-500">
                                <Grid3x3 size={20} />
                            </div>
                            <div>
                                <h3 class="text-base font-bold text-dark-gray">Échelle de la carte</h3>
                                <p class="text-xs text-stone-400">Configurez la grille pour dimensionner les jetons correctement · <strong class="text-stone-500">1 case = 1 m²</strong></p>
                            </div>
                        </div>

                        <!-- Scale Inputs -->
                        <div class="grid grid-cols-3 gap-4">
                            <div class="space-y-1.5">
                                <label for="pixels-per-cell" class="text-xs font-bold text-dark-gray">Pixels par case <span class="font-normal text-stone-400">(1 case = 1 m²)</span></label>
                                <input
                                    id="pixels-per-cell"
                                    type="number"
                                    min="10"
                                    max="500"
                                    step="1"
                                    bind:value={pixelsPerCell}
                                    class="w-full px-3 py-2 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all font-mono text-sm text-dark-gray"
                                />
                            </div>
                            <div class="space-y-1.5">
                                <label for="grid-offset-x" class="text-xs font-bold text-dark-gray">Décalage X (px)</label>
                                <input
                                    id="grid-offset-x"
                                    type="number"
                                    step="1"
                                    bind:value={gridOffsetX}
                                    class="w-full px-3 py-2 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all font-mono text-sm text-dark-gray"
                                />
                            </div>
                            <div class="space-y-1.5">
                                <label for="grid-offset-y" class="text-xs font-bold text-dark-gray">Décalage Y (px)</label>
                                <input
                                    id="grid-offset-y"
                                    type="number"
                                    step="1"
                                    bind:value={gridOffsetY}
                                    class="w-full px-3 py-2 bg-stone-50 border border-stone-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all font-mono text-sm text-dark-gray"
                                />
                            </div>
                        </div>

                        <!-- Pixel-per-cell slider -->
                        <div class="space-y-1.5">
                            <div class="flex justify-between items-center">
                                <span class="text-xs font-medium text-stone-400">Taille de la case</span>
                                <span class="text-xs font-mono font-bold text-blue-600">{pixelsPerCell}px</span>
                            </div>
                            <input
                                type="range"
                                min="10"
                                max="300"
                                step="1"
                                bind:value={pixelsPerCell}
                                class="w-full h-1.5 bg-stone-200 rounded-full appearance-none cursor-pointer accent-blue-500"
                            />
                            <div class="flex justify-between text-[10px] text-stone-300 font-mono">
                                <span>10</span>
                                <span>70 (défaut)</span>
                                <span>300</span>
                            </div>
                        </div>

                        <!-- Calibration Canvas -->
                        <div class="space-y-2">
                            <div class="flex items-center justify-between">
                                <span class="text-xs font-bold text-dark-gray">Calibration visuelle</span>
                                <span class="text-[10px] text-stone-400">Naviguez avec la molette • Zoomez avec Ctrl+Molette</span>
                            </div>
                            <div
                                bind:this={calibrationContainer}
                                bind:offsetWidth={calWidth}
                                bind:offsetHeight={calHeight}
                                class="relative w-full rounded-xl overflow-hidden border border-stone-200 bg-stone-900 select-none"
                                style="height: 320px; cursor: {calIsPanning ? 'grabbing' : 'grab'};"
                                onmousedown={handleCalMouseDown}
                                role="img"
                                aria-label="Calibration visuelle de la grille"
                            >
                                <canvas
                                    bind:this={calibrationCanvas}
                                    width={calWidth * (typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1)}
                                    height={calHeight * (typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1)}
                                    class="absolute inset-0 block"
                                    style="width: {calWidth}px; height: {calHeight}px;"
                                ></canvas>

                                <!-- Zoom Controls -->
                                <div class="absolute bottom-2 right-2 flex gap-1 bg-black/60 backdrop-blur-md p-1 rounded-lg border border-white/10 z-10">
                                    <button onclick={calZoomOut} class="p-1 text-stone-400 hover:text-white hover:bg-white/10 rounded transition-colors cursor-pointer" title="Zoom arrière">
                                        <ZoomOut size={14} />
                                    </button>
                                    <button onclick={calResetView} class="px-2 py-1 text-white hover:bg-white/10 rounded text-[10px] font-bold font-mono transition-colors cursor-pointer min-w-[40px] text-center">
                                        {Math.round(calZoom * 100)}%
                                    </button>
                                    <button onclick={calZoomIn} class="p-1 text-stone-400 hover:text-white hover:bg-white/10 rounded transition-colors cursor-pointer" title="Zoom avant">
                                        <ZoomIn size={14} />
                                    </button>
                                </div>
                            </div>
                            <p class="text-[10px] text-stone-400 leading-relaxed">
                                Ajustez les valeurs ci-dessus pour que la <strong class="text-blue-500">grille bleue</strong> corresponde aux cases de votre carte. 
                                Utilisez le décalage X/Y si la grille de la carte ne commence pas au coin haut-gauche.
                            </p>
                        </div>
                    </div>
                {/if}

                <!-- Form Actions -->
                <div class="pt-6 border-t border-stone-100 flex justify-end gap-3">
                    <a
                        href="/table/{gameId}/gm/settings?tab=boards"
                        class="px-5 py-2.5 bg-white border border-stone-200 hover:bg-stone-50 text-stone-600 rounded-xl font-medium transition-all text-sm"
                    >
                        Annuler
                    </a>
                    <button
                        onclick={handleSave}
                        disabled={saving || isUploadingFile || !boardName.trim() || !boardUrl.trim()}
                        class="px-5 py-2.5 bg-burnt-orange hover:bg-burnt-orange-hover text-white rounded-xl font-medium disabled:opacity-50 transition-all text-sm shadow-sm flex items-center justify-center min-w-[100px]"
                    >
                        {#if saving}
                            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                        {:else}
                            Enregistrer
                        {/if}
                    </button>
                </div>
            </div>
        {/if}
    </main>
</div>
