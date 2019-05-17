function renderStars(id, rating) {
    let result = ``;
    for (let i = 1.0; i <= 5.0; i++) {
        if (rating >= i) {
            result += `<span class="fa fa-star checked"></span>`;
        }else if (rating < i && rating > i - 1.0) {
            result += `<span class="fa fa-star-half-empty checked"></span>`;
        }else if (rating < i) {
            result += `<span class="fa fa-star"></span>`;
        }
    }
    result += `<span>(${rating})</span>`
    document.getElementById(id).innerHTML = result;
    return result
}

function renderStarsComment(id, rating) {
    let result = ``;
    for (let i = 1.0; i <= 5.0; i++) {
        if (rating >= i) {
            result += `<span class="fa fa-star checked"></span>`;
        }else if (rating < i && rating > i - 1.0) {
            result += `<span class="fa fa-star-half-empty checked"></span>`;
        }else if (rating < i) {
            result += `<span class="fa fa-star"></span>`;
        }
    }
    result += `<span>(${rating})</span>`
    document.getElementById(id).innerHTML = result;
    return result
}

// document.getElementsByName("rating").innerHTML = renderStars();