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
    // document.getElementById(id).innerHTML = result;
    return result
}

function sendMsg(receiver_id) {
    let url = '/send_message/user_id' + receiver_id;
    let text = document.getElementById('message');
    let msg = {
        'receiver_id': receiver_id,
        'text':text.value,
        'read':false,
        'date_send':new Date(),
    };

    if (text.value !== ''){
        console.log(msg);
        fetch(url,
            {
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                method: "POST",
                body: JSON.stringify(msg)
            })
        // .then(function(res){ console.log(res.body) })
            .catch(function(res){ console.log(res) })

        text.value = '';
        let modal = document.getElementById('modal-send-message');
        let instance = M.Modal.getInstance(modal);
        instance.close()
    }else {
        console.log('empty');
    }
}