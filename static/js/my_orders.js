function collapsible() {
    var elems = document.querySelectorAll('.collapsible');
    var instances = M.Collapsible.init(elems, {});
}

const select = document.getElementById('status');
const orders = document.getElementById('orders');

window.onload = function(){
    getOrders();
};

function getOrders() {
    fetch('/my_profile/sort_orders_by_'+select.value)
        .then(res => res.json())
        .then((json) => {
            orders.innerText = ``;
            if (select.value === select.options[0].value){
                renderAvailableOrders(json);
            }else if (select.value === select.options[1].value){
                renderPerformedOrders(json)
            }else if (select.value === select.options[2].value){
                renderDoneOrders(json)
            }
        })
        .catch(err => { throw err });
}

function renderAvailableOrders(availableOrders) {
    if (availableOrders === null){
        renderOrders(new NoOrders(select.value).create());
    }else {
        for (let i in availableOrders) {
            let order = new CardAvailableOrder(availableOrders[i]);
            renderOrders(order.create());
        }
    }
}

function renderPerformedOrders(performedOrders) {
    if (performedOrders === null){
        renderOrders(new NoOrders("in work").create());
    }else {
        for (let i in performedOrders) {
            let order = new CardPerformedOrder(performedOrders[i]);
            renderOrders(order.create());
        }
    }
}

function renderDoneOrders(doneOrders) {
    if (doneOrders === null){
        renderOrders(new NoOrders(select.value).create());
    }else {
        for (let i in doneOrders) {
            let order = new CardDoneOrder(doneOrders[i]);
            renderOrders(order.create());
            collapsible();
        }
    }
}

function renderOrders(child){
    orders.appendChild(child);
}

class CardDoneOrder{
    constructor(doneOrder) {
        this.doneOrder = doneOrder;
    }

    create(){
        let card = document.createElement("div");
        card.className = 'card order';

        card.appendChild(this.createContent());
        card.appendChild(this.createAction());

        return card;
    }

    createContent(){
        let card_content = new CardContent().create();
        let freelancer_comment = this.doneOrder.freelancer_comment.text;
        if (freelancer_comment === ""){
            freelancer_comment = "Volunteer did not comment";
        }
        card_content.innerHTML = `<span class="card-title">
                                     <a href="/orders/id${this.doneOrder.order.id}"> ${this.doneOrder.order.title} </a>
                                  </span>
                                  <ul class="collapsible">
                                    <li>
                                      <div class="collapsible-header card-inform">
                                          <span>Customer comment</span>
                                          <span> <a href="/customers/id"> ${this.doneOrder.order.customer.FirstName} ${this.doneOrder.order.customer.LastName} </a></span>
                                          <div class="rating" id="${this.doneOrder.customer_comment.id}">
                                                ${renderStarsComment(this.doneOrder.customer_comment.id, this.doneOrder.customer_comment.rait)}
                                          </div>
                                      </div>
                                      <div class="collapsible-body"><span>${this.doneOrder.customer_comment.text}</span></div>
                                    </li>
                                    <li>
                                      <div class="collapsible-header card-inform">
                                        <span>Volunteer comment</span>
                                          <span> <a href="freelancers/id${this.doneOrder.freelancer.ID}"> ${this.doneOrder.freelancer.FirstName} ${this.doneOrder.freelancer.LastName}</a></span>
                                          <div class="rating" id="${this.doneOrder.freelancer_comment.id}">
                                              ${renderStarsComment(this.doneOrder.freelancer_comment.id, this.doneOrder.freelancer_comment.rait)}
                                          </div>
                                      </div>
                                      <div class="collapsible-body"><span>${freelancer_comment}</span></div>
                                    </li>
                                  </ul>`;
        return card_content
    }

    createAction(){
        let card_action = new CardAction().create();
        card_action.classList.add('card-inform');
        card_action.innerHTML = `<span class="center">Complete: ${formatDate(new Date(this.doneOrder.date_complete))}</span>
                                <span class="center"> <a href="/freelancers/id${this.doneOrder.freelancer.ID}">${this.doneOrder.freelancer.FirstName} ${this.doneOrder.freelancer.LastName}</a></span>`;

        return card_action
    }
}

class CardPerformedOrder{
    constructor(availableOrder) {
        this.performedOrder = availableOrder;
    }

    create(){
        let card = document.createElement("div");
        card.className = 'card order';
        console.log(this.performedOrder);
        card.appendChild(this.createContent());
        card.appendChild(this.createAction());
        return card;
    }

    createContent(){
        let card_content = new CardContent().create();
        card_content.innerHTML = `<span class="card-title">
                                     <a href="/orders/id${this.performedOrder.order.id}"> ${this.performedOrder.order.title} </a>
                                  </span>
                                  <span>${this.performedOrder.order.content}</span>`;

        return card_content
    }

    createAction(){
        let card_action = new CardAction().create();
        card_action.innerHTML = `<p class="green-text">${this.performedOrder.order.status.Name}</p>
                                <span><a href="/freelancers/id${this.performedOrder.freelancer.id}">${this.performedOrder.freelancer.FirstName} ${this.performedOrder.freelancer.LastName}</a></span> 
                                <div class="action-block right">
                                    <a href="/orders/id${this.performedOrder.order.id}/make_done" title="Done"><i class="material-icons  light-blue-text darken-1 right">check_circle</i></a>
                                </div>`;

        return card_action
    }
}

class CardAvailableOrder{
    constructor(performedOrder) {
        this.order = performedOrder;
    }

    create(){
        let card = document.createElement("div");
        card.className = 'card order';

        card.appendChild(this.createContent());
        card.appendChild(this.createAction());

        return card;
    }

    createContent(){
        let card_content = new CardContent().create();
        let request_count = 0;
        if (this.order.freelancer_request !== null){request_count = this.order.freelancer_request.length}
        card_content.innerHTML = `<span class="card-title">
                                     <a href="/orders/id${this.order.id}"> ${this.order.title} </a>
                                  </span>
                                  <p class="bold">Requests ${request_count}</p>
                                  <span>${this.order.content}</span>`;

        return card_content
    }

    createAction(){
        let card_action = new CardAction().create();
        card_action.innerHTML = `<p class="green-text">${this.order.status.Name}</p>
                                <span>${formatDate(new Date(this.order.created_at))}</span>
                                <div class="action-block right">
                                    <a href="/my_profile/delete_order/id${this.order.id}" title="Delete"><i class="material-icons  light-blue-text darken-1 right">delete</i></a>
                                    <a href="/my_profile/edit_order/id${this.order.id}" title="Edit"><i class="material-icons  light-blue-text darken-1 right">edit</i></a>
                                </div>`;

        return card_action
    }
}

class CardContent{
    constructor(){}
    create(){
        let card_content = document.createElement("div");
        card_content.className = 'card-content';

        return card_content;
    }
}

class CardAction{
    constructor(){}
    create(){
        let card_action = document.createElement("div");
        card_action.className = 'card-action';

        return card_action;
    }
}

class NoOrders {
    constructor(statusName){
        this.text= 'No ' + statusName + ' orders';
    }
    create(){
        let no_orders = document.createElement('p');
        no_orders.className = 'center-align flow-text';
        no_orders.innerText = this.text;

        return no_orders;
    }
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