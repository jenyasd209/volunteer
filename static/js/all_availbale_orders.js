function deleteAvailableOrder(id) {
    let order_id = {'id':id};
    let url = "/moderator/delete_order_id"+id;
    console.log(JSON.stringify(order_id));
    fetch(url,
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(order_id)
        })
        .catch(function(res){ console.log(res) });
    document.getElementById("order-id"+id).remove();
    alert( "Order successful removed!");
}