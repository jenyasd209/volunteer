window.onload = function(){
    fetch('/my_profile/user_dialogs')
        .then(res => res.json())
        .then((json) => {
            renderInboxChat(json);
        })
        .catch(err => { throw err });
};

function renderInboxChat(dialogs) {
    let inbox_chat = document.getElementById('inbox_chat');
    let html = ``;
    if (dialogs === null){
        html += `<div class="valign-wrapper" style="height: 100%; justify-content: center">
                        <h5>No dialogs</h5>
                    </div>`
    }else {
        for (let i in dialogs){
            let indexLastMsg = dialogs[i]['messages'].length - 1;
            let last_msg = dialogs[i]['messages'][indexLastMsg];
            let read_status = "";
            if (!last_msg.read){
                if (last_msg.sender_id === dialogs[i].user_current.ID){
                    read_status = `<span class="grey-dot right"></span>`;
                }else {
                    read_status = `<span class="red-dot right"></span>`;
                }
            }
            html += `<div class="chat_list" onclick="viewDialog(${dialogs[i]['id']}, this)">
                     <div class="chat_people">
                         <div class="chat_img"><img class="avatar" src="${dialogs[i]['user_two']['Photo']}" alt="sunil"></div>
                         <div class="chat_ib">
                             <h5>${dialogs[i]['user_two']['FirstName']} ${dialogs[i]['user_two']['LastName']}
                                <span class="chat_date">${formatDate(new Date(dialogs[i]['messages'][indexLastMsg]['date_send']))}</span>   
                             </h5>
                             <p id="last-msg_id${dialogs[i]['id']}">${last_msg.text} ${read_status}</p>
                         </div>
                     </div>
                 </div>`
        }
    }
    addInnerHTML(inbox_chat, html);
}

function viewDialog(id, elem){
    let url = `/my_profile/dialog/id`+id;

    fetch(url)
        .then(res => res.json())
        .then((json) => {
            renderMsgHistory(json);
            setActive(elem);
            let msg_history = document.getElementById('msg_history');
            msg_history.scrollTop = msg_history.scrollHeight;
        })
        .catch(err => { throw err });
}

function setActive(elem) {
    // let chat_list = document.getElementsByClassName('chat_list');
    let current  = document.getElementsByClassName('active_chat');
    if (current.length !== 0){
        current[0].className = current[0].className.replace(" active_chat", "");
    }
    elem.classList.add('active_chat');
}

function renderMsgHistory(dialog){
    let input_field = `<div class="type_msg">
                           <div class="input_msg_write">
                               <input type="text" id="msg_text" class="write_msg" placeholder="Type a message" />
                               <button class="msg_send_btn" onclick="sendMsg(${dialog['user_current']['ID']}, ${dialog['user_two']['ID']}, ${dialog['id']})" type="button"><i class="fa fa-paper-plane-o" aria-hidden="true"></i></button>
                           </div>
                       </div>`;
    let msg_history = document.createElement("div");
    msg_history.id = 'msg_history';
    msg_history.className = "msg_history";
    // let msg_history = document.getElementById('msg_history');
    // msg_history.innerText = '';
    let messages = dialog['messages'];

    document.getElementById('mesgs').innerHTML = ``;

    for (let i in messages) {
        let msg_text = `<p>${messages[i]['text']}</p>`
        let msg_date_send = `<span class=\"time_date\">${formatDate(new Date(messages[i]['date_send']))}</span>`
        if (dialog['user_current']['ID'] === messages[i]['sender_id']){
            msg_history.innerHTML += `<div class="outgoing_msg">
                           <div class="sent_msg">
                               ${msg_text}
                               ${msg_date_send}
                           </div>
                       </div>`;
        } else{
            msg_history.innerHTML += `<div class="incoming_msg">
                            <div class="incoming_msg_img"> <img class="avatar" src="${dialog['user_two']['Photo']}" alt="sunil"> </div>
                            <div class="received_msg">
                                <div class="received_withd_msg">
                                    ${msg_text}
                                    ${msg_date_send}
                                </div>
                            </div>
                       </div>`;
        }
    };
    document.getElementById('mesgs').appendChild(msg_history);
    addInnerHTML(document.getElementById('mesgs'), input_field);
}

function formatDate(date) {
    let days = ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'];
    let months = ['January','February','March','April','May','June','July','August','September','October','November','December'];

    let formatDate = date.getHours() + ':' + date.getMinutes() + ' | ';
    if (date.getDate() === new Date().getDate()){
        formatDate += `Today`;
    }else if (date.getDate() === (new Date().setDate(new Date().getDate() - 1))){
        formatDate += `Yesterday`
    }else{
        formatDate += date.getDate() + '.' + date.getMonth() + '.' + date.getFullYear();
    }
    return formatDate
}

function addInnerHTML(elem, html) {
    elem.innerHTML += html;
}

function sendMsg(sender_id, receiver_id, dialog_id) {
    let url = '/send_message/user_id' + receiver_id;
    let text = document.getElementById('msg_text');

    let msg = {
        'sender_id':sender_id,
        'receiver_id': receiver_id,
        'dialog_id':dialog_id,
        'text':text.value,
        'read':false,
        'date_send':new Date(),
    };
    if (msg.text !== ""){
        let msg_text = `<p>${msg.text}</p>`;
        let msg_date_send = `<span class=\"time_date\">${formatDate(msg.date_send)}</span>`;
        let newMsg = `<div class="outgoing_msg">
                      <div class="sent_msg">
                          ${msg_text}
                          ${msg_date_send}
                      </div>
                  </div>`;
        let msg_history = document.getElementById('msg_history');

        fetch(url,
            {
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                method: "POST",
                body: JSON.stringify(msg)
            })
            .catch(function(res){ console.log(res) })

        text.value = '';
        msg_history.innerHTML += newMsg;
        msg_history.scrollTop = msg_history.scrollHeight;
        lastMessageRender(msg.text, dialog_id);
    }
}

function lastMessageRender(msg_text, dialog_id) {
    let last_msg = document.getElementById('last-msg_id'+dialog_id);
    last_msg.innerText = msg_text;
    last_msg.innerHTML += `<span class="grey-dot right"></span>`;
}