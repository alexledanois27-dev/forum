// Cada grupo permite «Général» O varias subcategorías. Los grupos se pueden
// combinar entre sí. Los IDs reales vienen del HTML generado por el backend.
(() => {
    const groups = Array.from(document.querySelectorAll(".post-category-group"));
    const status = document.querySelector(".category-selection-status");
    if (!groups.length || !status) return;

    function updateSelection() {
        let total = 0;
        groups.forEach((group) => {
            const checked = group.querySelectorAll('input[type="checkbox"]:checked:not(:disabled)');
            const count = group.querySelector(".category-selection-count");
            count.textContent = String(checked.length);
            count.hidden = checked.length === 0;
            group.classList.toggle("has-selection", checked.length > 0);
            total += checked.length;

            // Al seleccionar una subcategoría, asociamos también el grupo padre
            // si existe. Así la consulta actual por categoría encuentra el post
            // tanto en «Lore» como en «Races», sin tener que modificar Go.
            // No duplicamos ese ID cuando el checkbox «Général» ya lo envía.
            const parent = group.querySelector("[data-parent-category]");
            if (parent) {
                parent.name = "categories";
                parent.disabled = !group.querySelector("[data-subcategory]:checked:not(:disabled)");
            }
        });
        status.hidden = false;
        status.textContent = total === 0 ? "Aucune catégorie sélectionnée." :
            `${total} sélection${total > 1 ? "s" : ""} — vous pouvez choisir plusieurs thèmes.`;
    }

    // Los botones comparten una sola zona visible, pero ocultar un panel
    // no borra sus selecciones ni desactiva sus campos al enviar el formulario.
    groups.forEach((group) => {
        const button = group.querySelector(".category-group-button");
        const panel = group.querySelector(".post-category-panel");
        panel.hidden = true;
        button.setAttribute("aria-expanded", "false");
        button.addEventListener("click", () => {
            const shouldOpen = panel.hidden;
            groups.forEach((other) => {
                other.querySelector(".post-category-panel").hidden = true;
                other.querySelector(".category-group-button").setAttribute("aria-expanded", "false");
            });
            panel.hidden = !shouldOpen;
            button.setAttribute("aria-expanded", String(shouldOpen));
        });
        panel.addEventListener("keydown", (event) => {
            if (event.key === "Escape") {
                event.preventDefault();
                panel.hidden = true;
                button.setAttribute("aria-expanded", "false");
                button.focus();
            }
        });
        group.addEventListener("change", (event) => {
            const input = event.target;
            if (input.matches("input[type=checkbox]") && input.checked) {
                // Le choix le plus récent remplace l'autre niveau de précision.
                const others = input.hasAttribute("data-general") ?
                    group.querySelectorAll("[data-subcategory]") :
                    group.querySelectorAll("[data-general]");
                others.forEach((other) => { other.checked = false; });
            }
            updateSelection();
        });
    });
    // Restaure aussi les compteurs après Retour navigateur ou reset du formulaire.
    window.addEventListener("pageshow", updateSelection);
    const form = groups[0].closest("form");
    form.addEventListener("reset", () => requestAnimationFrame(updateSelection));
    updateSelection();
})();
