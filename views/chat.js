//Setting Nickname
$('#nickname').text(localStorage.getItem("nickname"));
localStorage.removeItem("nickname");

//Setting Room code
$('#roomCode').text(localStorage.getItem("code"));
localStorage.removeItem("code");

//Styling chatBox
let totalHeight = $(window).height();
$('#chatBox').css('height', totalHeight - 200);

//Grapping essensial data
const nickname = $('#nickname').text();
const roomCode = $('#roomCode').text();

const socket = io.connect(window.Location.href);

//Joining the room
socket.on('initialConnection', data=>{
  console.log("Connected!");
  socket.emit('initialConnection', roomCode)
});

//Recieving a message
socket.on('chatMessage', data=>{
  let info = data.split("$");
  $("#chatBox").append('<br><p><strong>' + info[0] + ': </strong>'+ info[2] +'</p>');
  $('#typing').text('');
});


//Sending message by button or Enter key.
$('#sendMessage').on('click', ()=>{send()});
$('#messageBox').keypress(e=>{ if(e.which == 13) {send()}});

//Sending message function
function send(){
  socket.emit('chatMessage', nickname + "$" + roomCode + "$" + $('#messageBox').val());
  $('#messageBox').val("");
}


//Tracking typing
$('#messageBox').keypress(e=>{
    socket.emit('typing', nickname + "$" + roomCode);
});

//Someone else is typing
socket.on('typing', data=>{
  let info = data.split("$");
  $('#typing').text(info[0] + ' is typing...');
  setTimeout(()=>{
    $('#typing').text('');
  },5000);
})