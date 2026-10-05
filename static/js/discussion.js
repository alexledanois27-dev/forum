// El título lleva al detalle sin abrir el formulario. El icono de comentario
// añade #reply a la URL para abrirlo y colocar el cursor en el campo de texto.
(() => {
    function openReply() {
        if (window.location.hash !== "#reply") return;
        const reply = document.getElementById("reply");
        if (!reply) return;
        if (reply.tagName === "DETAILS") reply.open = true;
        // Si no hay sesión, enfocamos la caja que invita a conectarse.
        const target = reply.querySelector("textarea") || reply;
        target.focus({ preventScroll: true });
        reply.scrollIntoView({ block: "start" });
    }
    window.addEventListener("hashchange", openReply);
    openReply();
})();
