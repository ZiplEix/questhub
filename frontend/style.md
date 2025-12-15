# Charte Graphique - QuestHub

Ce document recense les règles de style et les tokens de design utilisés sur le projet QuestHub.

## 🎨 Palette de Couleurs

L'identité visuelle repose sur des tons chauds et naturels évoquant le papier et l'aventure classique, contrastés par un gris bleuté moderne.

| Nom | Variable CSS / Tailwind | Code Hex | Usage |
| :--- | :--- | :--- | :--- |
| **Cream** | `--color-cream` / `bg-cream` | `#F9F7F2` | **Fond principal** de l'application. Remplace le blanc pur pour réduire la fatigue visuelle. |
| **Burnt Orange** | `--color-burnt-orange` / `text-burnt-orange` | `#E07A5F` | **Couleur Primaire**. Boutons d'action, liens importants, accents vifs. |
| **Dark Gray** | `--color-dark-gray` / `text-dark-gray` | `#3D405B` | **Texte principal**, titres, et fonds sombres (navbar, cartes actives). |
| **Mustard Yellow** | `--color-mustard-yellow` / `bg-mustard-yellow` | `#F2CC8F` | **Secondaire**. Éléments décoratifs, illustrations, états d'avertissement doux. |
| **Soft Shadow** | `--color-soft-shadow` | `rgba(0, 0, 0, 0.05)` | Ombres légères pour la profondeur. |

### Nuances Utilitaires
*   **Blanc (`#FFFFFF`)** : Utilisé pour les cartes (`cards`) sur fond Crème.
*   **Stone 50/100/200** : Utilisé pour les bordures subtiles et les fonds de conteneurs secondaires (ex: zones de saisie).

---

## ✒️ Typographie

| Rôle | Police | Classes Tailwind | Usage |
| :--- | :--- | :--- | :--- |
| **Titres** | **Outfit** | `font-display` | Tous les niveaux de titres (`h1` à `h6`). Apporte un côté moderne et "Roleplay". |
| **Corps** | **Inter** | `font-sans` | Texte courant, paragraphes, boutons. Optimisé pour la lisibilité UI. |

---

## 📐 Formes & Espacements

Le design utilise des formes très arrondies pour un rendu amical et organique.

### Border Radius
*   **XL (`1.5rem` / `24px`)** : Standard pour les **boutons**, **inputs**, et petits conteneurs.
    *   *Utility:* `rounded-xl`
*   **2XL (`2rem` / `32px`)** : Standard pour les **cartes principales**, modales, et gros blocs de contenu.
    *   *Utility:* `rounded-2xl`
*   **Full** : Pour les badges et les pilules de navigation.
    *   *Utility:* `rounded-full`

### Ombres
*   Utilisation d'ombres douces (`shadow-sm`, `shadow-md`) pour détacher les cartes du fond Crème.
*   Éviter les ombres dures ou noires pures.

---

## 🧩 Composants UI Standards

### Boutons
1.  **Primaire** : Fond `bg-burnt-orange`, Texte `text-white`, `rounded-xl` ou `rounded-2xl`, Shadow `shadow-md`.
    *   *Hover* : Légère translation verticale (`-translate-y-0.5`) et opacité réduite.
2.  **Secondaire / Outline** : Fond `bg-white` ou transparent, Bordure `border-stone-200`, Texte `text-dark-gray`.
    *   *Hover* : `bg-stone-50`, Texte `text-burnt-orange`, Bordure `border-burnt-orange/30`.
3.  **Ghost / Tertiaire** : Pas de fond, Texte `text-stone-500`.
    *   *Hover* : `text-dark-gray`.

### Cartes (Cards)
*   Fond `bg-white`.
*   Bordure fine `border border-stone-100`.
*   Radius `rounded-2xl`.
*   Ombre légère `shadow-sm`.

### Inputs / Formulaires
*   Fond `bg-stone-50` (pour se détacher du blanc des cartes).
*   Bordure `border-stone-200`.
*   Radius `rounded-xl`.
*   Focus : Ring `burnt-orange`.

### Navigation (Tabs)
*   **Segmented Control** : Conteneur `bg-stone-100` `rounded-xl` avec padding `p-1`. Élément actif `bg-white` `shadow-sm` `rounded-[radius-interne]`.

---

## 🔮 Iconographie

*   **Librairie** : [Lucide Svelte](https://lucide.dev/)
*   **Style** : Traits (Stroke) `2px` ou `1.5px` selon la taille.
*   **Cohérence** : Toujours utiliser les composants Lucide (ex: `<Search />`, `<Dice5 />`) plutôt que des SVG inline hardcodés.
