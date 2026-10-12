// <details> abre y cierra cada categoría de forma nativa al hacer clic
// o al pulsar Enter/Espacio sobre su título. No usamos eventos de hover.
(() => {
    const groups = document.querySelectorAll(".category-group");
    groups.forEach((group) => {
        // El atributo name del HTML permite un solo grupo abierto.
        // Este evento mantiene el comportamiento en navegadores más antiguos.
        group.addEventListener("toggle", () => {
            if (!group.open) return;
            groups.forEach((other) => {
                if (other !== group) other.open = false;
            });
        });
        // Escape cierra el grupo y devuelve el foco a su título.
        group.addEventListener("keydown", (event) => {
            if (event.key === "Escape" && group.open) {
                event.preventDefault();
                event.stopPropagation();
                group.open = false;
                group.querySelector("summary").focus();
            }
        });
    });
})();
