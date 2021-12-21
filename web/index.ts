type Order = {
    order_uid: string;
    entry: string;
    total_price: number;
    customer_id: string;
    track_number: string;
    delivery_service: string;
}
  
async function get() {
    const order = (<HTMLInputElement>document.getElementById("order_uid")).value;
    
    const response = await window.fetch('http://127.0.0.1:3000', {
        method: 'GET',
        headers: {
          'order_uid': order,
          'Access-Control-Request-Origin': '*',
          'Access-Control-Request-Method': '*',
        }
    })
    
    let textResponse = await response.text()
    let nodeTextString = document.createTextNode(textResponse)
    document.body.append(nodeTextString)
}
