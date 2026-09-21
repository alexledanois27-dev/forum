// Ordenamos los posts que ya están en la página, sin llamar al backend.
(() => {
    const select = document.getElementById("feed-sort");
    const list = document.querySelector(".post-list");
    if (!select || !list) return;

    // Los atributos data-* contienen la fecha Unix y los likes entregados por Go.
    const cards = Array.from(list.querySelectorAll(".post-card"));
    select.disabled = cards.length === 0;

    const requestedSort = new URLSearchParams(window.location.search).get("sort");
    const popularView = requestedSort === "popular";
    let sortMode = requestedSort === "oldest" ? "oldest" : "newest";
    select.value = sortMode;

    function sortPosts() {
        cards.sort((a, b) => {
            const dateDifference = Number(b.dataset.created) - Number(a.dataset.created);
            const chronologicalOrder = sortMode === "oldest" ? -dateDifference : dateDifference;
            if (popularView) {
                // Conservamos la popularidad y combinamos el orden elegido:
                // primero los likes, después la fecha en caso de empate.
                return Number(b.dataset.likes) - Number(a.dataset.likes) || chronologicalOrder;
            }
            return chronologicalOrder;
        });

        // Movemos los elementos originales: no recreamos su HTML y conservamos
        // los comentarios escritos, los eventos y el estado de "Lire la suite".
        cards.forEach((card) => list.append(card));
    }

    // La vista y el orden cronológico son independientes: el desplegable
    // permanece disponible también al entrar desde «Populaires».
    if (popularView) {
        const heading = document.querySelector(".feed-sort-bar h2");
        if (heading) heading.textContent = "Chroniques populaires";
    }
    select.addEventListener("change", () => {
        sortMode = select.value;
        sortPosts();
    });
    sortPosts();
})();
