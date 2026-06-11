<script lang="ts">
    import { onMount } from "svelte";
    import { fetchMediaLibrary, uploadToMediaLibrary, type MediaAsset } from "$lib/api";
    import { X, Loader2, Camera, RefreshCw } from "lucide-svelte";

    let { 
        isOpen = false, 
        onSelect, 
        onClose 
    } = $props<{
        isOpen: boolean;
        onSelect: (url: string) => void;
        onClose: () => void;
    }>();

    let mediaAssets = $state<MediaAsset[]>([]);
    let loading = $state(false);
    let uploading = $state(false);
    let error = $state("");
    let fileInput = $state<HTMLInputElement | null>(null);

    onMount(async () => {
        if (isOpen) {
            await loadMedia();
        }
    });

    // Reload when modal opens
    $effect(() => {
        if (isOpen) {
            loadMedia();
        }
    });

    async function loadMedia() {
        try {
            loading = true;
            error = "";
            mediaAssets = await fetchMediaLibrary();
        } catch (e: any) {
            console.error("Failed to load media assets", e);
            error = "Erreur lors du chargement de la médiathèque.";
        } finally {
            loading = false;
        }
    }

    async function handleUpload(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        if (!file.type.startsWith('image/')) {
            error = "Veuillez sélectionner un fichier image.";
            return;
        }

        if (file.size > 10 * 1024 * 1024) { // 10MB limit
            error = "L'image ne doit pas dépasser 10 Mo.";
            return;
        }

        uploading = true;
        error = "";

        try {
            const asset = await uploadToMediaLibrary(file);
            onSelect(asset.url);
        } catch (e: any) {
            console.error("Failed to upload image", e);
            error = e.message || "Erreur lors du téléversement.";
        } finally {
            uploading = false;
            if (input) input.value = "";
        }
    }
</script>

{#if isOpen}
    <div class="fixed inset-0 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4 z-[999] animate-in fade-in duration-200">
        <div class="bg-white rounded-3xl p-6 max-w-3xl w-full border border-stone-200 shadow-2xl flex flex-col max-h-[85vh]">
            <!-- Header -->
            <div class="flex justify-between items-center pb-4 border-b border-stone-100">
                <div>
                    <h3 class="font-display font-black text-xl text-dark-gray">Sélectionner depuis la Médiathèque</h3>
                    <p class="text-xs text-stone-400">Choisissez une image existante ou téléversez-en une nouvelle.</p>
                </div>
                <button 
                    onclick={onClose}
                    class="p-1.5 text-stone-400 hover:text-stone-600 hover:bg-stone-100 rounded-full transition-all cursor-pointer"
                >
                    <X size={18} />
                </button>
            </div>

            <!-- Content Area (Scrollable) -->
            <div class="grow overflow-y-auto py-6 space-y-4 min-h-[300px]">
                {#if error}
                    <div class="bg-red-50 text-red-600 px-4 py-3 rounded-xl text-xs border border-red-100">
                        {error}
                    </div>
                {/if}

                {#if loading}
                    <div class="flex items-center justify-center py-20">
                        <Loader2 class="animate-spin text-burnt-orange" size={32} />
                    </div>
                {:else if mediaAssets.length > 0}
                    <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                        <!-- Direct upload card -->
                        <button
                            onclick={() => fileInput?.click()}
                            disabled={uploading}
                            class="bg-stone-50/50 hover:bg-stone-50 rounded-2xl border border-dashed border-stone-200 hover:border-burnt-orange/40 p-4 transition-all flex flex-col items-center justify-center gap-2 aspect-square cursor-pointer disabled:opacity-50"
                        >
                            {#if uploading}
                                <Loader2 size={24} class="animate-spin text-burnt-orange" />
                                <span class="text-xs font-bold text-stone-500">Téléversement...</span>
                            {:else}
                                <Camera size={24} class="text-stone-400" />
                                <span class="text-xs font-bold text-stone-500">Ajouter une image</span>
                            {/if}
                        </button>

                        {#each mediaAssets as asset}
                            <button
                                onclick={() => onSelect(asset.url)}
                                class="bg-white rounded-2xl border border-stone-150 overflow-hidden shadow-3xs hover:shadow-md hover:border-burnt-orange/50 transition-all flex flex-col justify-between aspect-square text-left cursor-pointer group relative"
                            >
                                <img
                                    src={asset.url}
                                    alt={asset.name}
                                    class="w-full h-full object-cover group-hover:scale-102 transition-all duration-300"
                                />
                                <div class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 to-transparent p-3 pt-6 text-white opacity-0 group-hover:opacity-100 transition-opacity">
                                    <p class="text-[10px] font-bold truncate">{asset.name}</p>
                                    <p class="text-[8px] text-stone-300 font-mono mt-0.5">{Math.round(asset.size / 1024)} Ko</p>
                                </div>
                            </button>
                        {/each}
                    </div>
                {:else}
                    <div class="text-center py-16 bg-stone-50/50 rounded-2xl border border-stone-100 flex flex-col items-center justify-center gap-4">
                        <div class="w-12 h-12 bg-white rounded-full flex items-center justify-center text-stone-300 shadow-sm border border-stone-100">
                            <Camera size={24} />
                        </div>
                        <div class="space-y-1">
                            <h4 class="font-bold text-dark-gray text-sm">Votre médiathèque est vide</h4>
                            <p class="text-stone-400 text-xs max-w-xs mx-auto">
                                Vous n'avez pas encore d'images stockées.
                            </p>
                        </div>
                        <button
                            onclick={() => fileInput?.click()}
                            disabled={uploading}
                            class="px-4 py-2 bg-burnt-orange hover:bg-opacity-95 text-white font-bold rounded-xl text-xs transition-all flex items-center gap-2 cursor-pointer shadow-sm"
                        >
                            {#if uploading}
                                <Loader2 size={14} class="animate-spin" />
                                Téléversement...
                            {:else}
                                Importez votre première image
                            {/if}
                        </button>
                    </div>
                {/if}
            </div>

            <!-- Hidden File Input -->
            <input
                type="file"
                accept="image/*"
                bind:this={fileInput}
                class="hidden"
                onchange={handleUpload}
            />

            <!-- Footer / Cancel Button -->
            <div class="flex justify-end pt-4 border-t border-stone-100">
                <button
                    type="button"
                    onclick={onClose}
                    class="px-4 py-2 border border-stone-200 hover:bg-stone-50 text-stone-600 font-bold rounded-xl text-xs transition-colors cursor-pointer"
                >
                    Annuler
                </button>
            </div>
        </div>
    </div>
{/if}
