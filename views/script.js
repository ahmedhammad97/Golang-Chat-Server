$('#create').on('click', ()=>{
    if(emptyNickname()){return;}
  
    //REQUESTING A ROOM FROM THE SERVER
    $.ajax({
      type: 'POST',
      data: JSON.stringify({'nickname' :  $('#nickname').val()}),
      url: '/create',
      success: function(data) {
        startChat(data, $('#nickname').val());
      },
      error: function(data){
        $('#message').text('Something went wrong, check your connection!');
      }
    });
  
  });
  
  
  $('#join').on('click', ()=>{
    if(emptyNickname()){return;}
  
    //ASKING FOR PERMISSION TO JOIN THE ROOM
    $.ajax({
      type: 'POST',
      data: JSON.stringify({'nickname' : $('#nickname').val(), 'roomCode' :  $('#roomCode').val()}),
      url: '/join',
      success: function(data) {
        if(data === "Not available"){
          $('#message').text("No such a room code");
        }
        else if (data === "Name exists") {
            $('#message').text("Nickname is already taken");
        }
        else{
          startChat($('#roomCode').val(), $('#nickname').val());
        }
      },
      error: function(data){
        $('#message').text('Something went wrong, check your connection!');
      }
    });
  
  });
  
  
  
  // HELPER FUNCTIONS
  function emptyNickname(){
    if($('#nickname').val().trim() == ""){
      $('#message').text('A nickname is required.');
      return true;
    }
    return false;
  };
  
  function startChat(code, nickname){
    localStorage.setItem("nickname", nickname);
    window.location.replace(`${window.location.href}room/${code}`);
  }