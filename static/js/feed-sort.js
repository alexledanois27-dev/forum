// Ordenamos los posts que ya están en la página, sin llamar al backend.
(() => {
    const select = document.getElementById("feed-sort");
    const list = document.querySelector(".post-list");
    if (!select || !list) return;

    // Los atributos data-* contienen la fecha Unix y los likes entregados por Go.
    const cards = Array.from(list.querySelectorAll(".post-card"));
    select.disabled = cards.length === 0;

    function sortPosts() {
        cards.sort((a, b) => {
            const dateDifference = Number(b.dataset.created) - Number(a.dataset.created);
            if (select.value === "oldest") return -dateDifference;
            if (select.value === "popular") {
                // Más likes primero. En caso de empate, primero el más reciente.
                return Number(b.dataset.likes) - Number(a.dataset.likes) || dateDifference;
            }
            return dateDifference;
        });

        // Movemos los elementos originales: no recreamos su HTML y conservamos
        // los comentarios escritos, los eventos y el estado de "Lire la suite".
        cards.forEach((card) => list.append(card));
    }

    select.addEventListener("change", sortPosts);
    sortPosts();
})();
