// El menú permanece en su lugar original. Solo lo movemos al panel
// mientras está abierto en móvil; al cerrarlo, lo devolvemos al mismo lugar.
(() => {
    const toggle = document.querySelector(".menu-toggle");
    const sidebar = document.getElementById("forum-navigation");
    const dialog = document.getElementById("mobile-menu");
    // Login y registro usan otro layout y no tienen menú lateral.
    if (!toggle || !sidebar || !dialog || typeof dialog.showModal !== "function") return;

    const closeButton = dialog.querySelector(".menu-close");
    const container = dialog.querySelector(".mobile-menu-content");
    const mobile = window.matchMedia("(max-width: 640px)");
    // Este marcador conserva el lugar exacto del menú en escritorio.
    const placeholder = document.createComment("Posición del menú de escritorio");
    sidebar.before(placeholder);

    function closeMenu() {
        if (dialog.open) dialog.close();
        placeholder.after(sidebar);
        toggle.setAttribute("aria-expanded", "false");
        document.body.classList.remove("menu-open");
    }

    function adaptMenu() {
        const focusWasInside = sidebar.contains(document.activeElement);
        closeMenu();
        if (mobile.matches) {
            toggle.hidden = false;
            if (focusWasInside) toggle.focus();
        } else {
            placeholder.after(sidebar);
            // Evitamos dejar el foco en el botón hamburguesa cuando se oculta.
            if (document.activeElement === toggle || dialog.contains(document.activeElement)) {
                sidebar.querySelector("a").focus();
            }
            toggle.hidden = true;
        }
    }

    toggle.addEventListener("click", () => {
        // showModal bloquea la interacción con el fondo y mantiene el foco
        // del teclado dentro del diálogo. Escape lo cierra de forma nativa.
        if (!mobile.matches) return;
        container.append(sidebar);
        dialog.showModal();
        toggle.setAttribute("aria-expanded", "true");
        document.body.classList.add("menu-open");
        closeButton.focus();
    });
    closeButton.addEventListener("click", closeMenu);
    dialog.addEventListener("cancel", (event) => {
        event.preventDefault();
        closeMenu();
    });
    dialog.addEventListener("close", () => {
        // El navegador devuelve el foco al botón que abrió el diálogo.
        // No alteramos el estado si ya se ha vuelto a abrir.
        if (!dialog.open) closeMenu();
    });

    // El fondo de un dialog recibe eventos como si fueran del propio dialog.
    // Comprobamos las coordenadas para distinguirlo de un clic dentro del panel.
    function outsidePanel(event) {
        const rect = dialog.getBoundingClientRect();
        return event.clientX < rect.left || event.clientX > rect.right ||
            event.clientY < rect.top || event.clientY > rect.bottom;
    }
    let startedOutside = false;
    dialog.addEventListener("pointerdown", (event) => {
        startedOutside = event.target === dialog && outsidePanel(event);
    });
    dialog.addEventListener("click", (event) => {
        if (startedOutside && event.target === dialog && outsidePanel(event)) closeMenu();
        startedOutside = false;
        if (event.target.closest("a[href]")) closeMenu();
    });

    // El cambio se hace al cruzar el mismo umbral definido en el CSS.
    mobile.addEventListener("change", adaptMenu);
    document.documentElement.classList.add("navigation-ready");
    adaptMenu();
})();
