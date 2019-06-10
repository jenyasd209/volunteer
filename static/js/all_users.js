function deleteUser(id) {
    let user_id = {'id':id};
    let url = "/moderator/delete_user_id"+id;
    console.log(JSON.stringify(user_id));
    fetch(url,
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(user_id)
        })
        .catch(function(res){ console.log(res) });
    document.getElementById("user-id"+id).remove();
    alert( "User successful removed!");
}