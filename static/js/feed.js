// Este archivo se carga con "defer": el navegador espera a tener el HTML
// preparado antes de ejecutarlo. No necesita ninguna librería.
(() => {
    // querySelectorAll encuentra todos los botones de lectura de la homepage.
    // Cada post conserva su propio estado: abrir uno no abre los demás.
    const posts = Array.from(document.querySelectorAll(".post-read-more"), (button) => {
        const content = document.getElementById(button.getAttribute("aria-controls"));
        let expanded = false;

        function update() {
            // Aplicamos temporalmente el límite CSS para medir seis líneas,
            // incluso cuando el usuario tiene el post desplegado.
            content.classList.add("is-collapsed");

            // scrollHeight es la altura total del texto; clientHeight es la
            // altura visible. Toleramos 1 píxel por los redondeos del navegador.
            const overflows = content.scrollHeight > content.clientHeight + 1;

            // Si ahora cabe todo (por ejemplo, al ampliar la ventana), ya no
            // necesitamos el botón ni conservar el estado de lectura ampliada.
            if (!overflows) expanded = false;
            button.hidden = !overflows;
            content.classList.toggle("is-collapsed", overflows && !expanded);
            button.textContent = expanded ? "Réduire" : "Lire la suite";

            // aria-expanded comunica el estado a los lectores de pantalla.
            button.setAttribute("aria-expanded", String(expanded));
        }

        // Un botón HTML ya responde al clic, a Enter y a la barra espaciadora.
        button.addEventListener("click", () => {
            expanded = !expanded;
            update();
        });

        update();
        return { content, update };
    });

    // Las fuentes pueden llegar después del HTML y cambiar el número de líneas.
    if (document.fonts) {
        document.fonts.ready.then(() => posts.forEach((post) => post.update()));
    }

    // Observamos el ancho del texto para adaptarnos al móvil o a una ventana
    // redimensionada. Ignoramos cambios solo de altura para no crear un ciclo
    // de mediciones al pulsar "Lire la suite" o "Réduire".
    if ("ResizeObserver" in window) {
        posts.forEach(({ content, update }) => {
            let previousWidth = content.getBoundingClientRect().width;
            const observer = new ResizeObserver(([entry]) => {
                const width = entry.contentRect.width;
                if (width !== previousWidth) {
                    previousWidth = width;
                    update();
                }
            });
            observer.observe(content);
        });
    } else {
        // Alternativa para navegadores que no disponen de ResizeObserver.
        window.addEventListener("resize", () => posts.forEach((post) => post.update()));
    }
})();
