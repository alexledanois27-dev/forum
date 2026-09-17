// <details> abre y cierra el menú de forma nativa, incluso sin JavaScript.
// Añadimos el cierre al pulsar fuera, al salir con Tab o al pulsar Escape.
(() => {
    const menu = document.querySelector(".account-menu");
    if (!menu) return;
    const trigger = menu.querySelector("summary");

    document.addEventListener("click", (event) => {
        if (!menu.contains(event.target)) menu.open = false;
    });
    document.addEventListener("focusin", (event) => {
        if (!menu.contains(event.target)) menu.open = false;
    });
    menu.addEventListener("keydown", (event) => {
        if (event.key === "Escape" && menu.open) {
            event.preventDefault();
            menu.open = false;
            trigger.focus();
        }
    });
})();
