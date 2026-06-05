import {
    Sword,
    Wand2,
    AlertCircle,
    Package,
    Zap,
    Heart,
    Scroll,
    Key,
    Map,
    Compass,
    Backpack,
    Gem,
    Trophy,
    Wand,
    Axe,
    Shield,
    Wind,
    Cloud,
    Star,
    Moon,
    Sun,
    Book,
    Coins,
    Skull,
    Eye,
    HelpCircle,
    Crosshair,
    Droplets,
    Flame,
} from "lucide-svelte";
import type { ComponentType } from "svelte";
import type { Icon as LucideIcon } from "lucide-svelte";

export interface IconOption {
    name: string;
    label: string;
    category: string;
    component: ComponentType<LucideIcon>;
}

// Liste des icônes disponibles avec métadonnées
export const AVAILABLE_ICONS: IconOption[] = [
    // Armes
    { name: "sword", label: "Épée", category: "Armes", component: Sword },
    { name: "axe", label: "Hache", category: "Armes", component: Axe },
    { name: "bow", label: "Arc", category: "Armes", component: Crosshair },
    { name: "wand", label: "Baguette", category: "Armes", component: Wand2 },
    { name: "staff", label: "Bâton", category: "Armes", component: Scroll },
    
    // Armures et protection
    { name: "shield", label: "Bouclier", category: "Protection", component: Shield },
    
    // Objets magiques
    { name: "gem", label: "Gemme", category: "Magique", component: Gem },
    { name: "spell", label: "Sort", category: "Magique", component: Wand },
    { name: "scroll", label: "Parchemin", category: "Magique", component: Scroll },
    { name: "book", label: "Livre", category: "Magique", component: Book },
    { name: "key", label: "Clé", category: "Magique", component: Key },
    
    // Utilitaire
    { name: "backpack", label: "Sac à dos", category: "Utilitaire", component: Backpack },
    { name: "map", label: "Carte", category: "Utilitaire", component: Map },
    { name: "compass", label: "Boussole", category: "Utilitaire", component: Compass },
    { name: "torch", label: "Torche", category: "Utilitaire", component: Flame },
    { name: "coin", label: "Pièce", category: "Utilitaire", component: Coins },
    { name: "treasure", label: "Trésor", category: "Utilitaire", component: Trophy },
    
    // Éléments
    { name: "water", label: "Eau", category: "Éléments", component: Droplets },
    { name: "fire", label: "Feu", category: "Éléments", component: Flame },
    { name: "air", label: "Air", category: "Éléments", component: Wind },
    { name: "earth", label: "Terre", category: "Éléments", component: Cloud },
    
    // Divers
    { name: "health", label: "Santé", category: "Divers", component: Heart },
    { name: "mana", label: "Mana", category: "Divers", component: Zap },
    { name: "star", label: "Étoile", category: "Divers", component: Star },
    { name: "moon", label: "Lune", category: "Divers", component: Moon },
    { name: "sun", label: "Soleil", category: "Divers", component: Sun },
    { name: "eye", label: "Œil", category: "Divers", component: Eye },
    { name: "skull", label: "Crâne", category: "Divers", component: Skull },
    { name: "question", label: "Question", category: "Divers", component: HelpCircle },
    { name: "package", label: "Paquet", category: "Divers", component: Package },
];

// Mapping pour la compatibilité avec getIconComponent
const ICON_MAP: Record<string, ComponentType<LucideIcon>> = AVAILABLE_ICONS.reduce(
    (acc, icon) => {
        acc[icon.name] = icon.component;
        return acc;
    },
    {} as Record<string, ComponentType<LucideIcon>>
);

/**
 * Récupère l'icône lucide-svelte correspondant au nom donné
 * @param iconName - Nom de l'icône (case-insensitive)
 * @returns Composant lucide-svelte ou Package si non trouvé
 */
export function getIconComponent(iconName?: string): ComponentType<LucideIcon> {
    if (!iconName) {
        return Package;
    }
    
    const normalizedName = iconName.toLowerCase().trim();
    return ICON_MAP[normalizedName] || Package;
}
