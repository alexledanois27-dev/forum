// Un mismo script sirve para el formulario de post y todos los comentarios.
// data-counter indica el ID del contador asociado a cada campo de texto.
(() => {
    const numberFormat = new Intl.NumberFormat("fr-FR");

    document.querySelectorAll("textarea[data-counter]").forEach((field) => {
        const counter = document.getElementById(field.dataset.counter);
        const limit = field.maxLength;

        function updateCounter() {
            // maxlength limita la escritura y el pegado en el navegador.
            // Usamos .length, que cuenta en unidades UTF-16 igual que maxlength:
            // algunos emojis ocupan dos unidades. No es un contador de bytes.
            const remaining = Math.max(0, limit - field.value.length);
            counter.textContent = `${numberFormat.format(remaining)} / ${numberFormat.format(limit)}`;
            counter.hidden = false;
        }

        // "input" se dispara al escribir, borrar, pegar o cortar texto.
        field.addEventListener("input", updateCounter);
        // pageshow cubre los valores recuperados al volver con el navegador.
        window.addEventListener("pageshow", updateCounter);
        if (field.form) {
            field.form.addEventListener("reset", () => {
                // Esperamos a que el navegador haya vaciado el formulario.
                requestAnimationFrame(updateCounter);
            });
        }
        updateCounter();
    });
})();
