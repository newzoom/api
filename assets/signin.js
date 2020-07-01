$(document).ajaxStart(function () {
  return NProgress.start();
});

$(document).ajaxStop(function () {
  return NProgress.done();
});

$(document).ready(function () {
  function signInCallback(authResult) {
    if (authResult["code"]) {
      $("#customBtn")[0].style.display = "none";
      $.ajax({
        type: "POST",
        url: "http://" + document.location.host + "/auth",
        headers: {
          "X-Requested-With": "XMLHttpRequest",
        },
        contentType: "application/octet-stream; charset=utf-8",
        processData: false,
        data: authResult["code"],
        success: function ({ data }) {
          console.log("sign in successfully");
          console.log(data);
          var { id, email, name, avatar, access_token } = data;
          document.cookie = `access_token=${access_token};`;
          document.cookie = `uid=${id};`;
          document.cookie = `email=${email};`;
          document.cookie = `avatar=${avatar};`;
          document.cookie = `name=${name};`;
          $("#name")[0].innerText = "Signed in: " + name;
        },
      });
    } else {
      alert("Error occurred, please comeback later");
    }
  }

  $("#customBtn").click(function () {
    auth2.grantOfflineAccess().then(signInCallback);
  });
});

function start() {
  gapi.load("auth2", function () {
    auth2 = gapi.auth2.init({
      client_id:
        "44490805046-2t5kjoq5s5jqh6isvm958321fqkspc1j.apps.googleusercontent.com",
    });
  });
}

function signOut() {
  var auth2 = gapi.auth2.getAuthInstance();
  auth2.signOut().then(function () {
    console.log("User signed out.");
    document.cookie =
      "access_token=; uid=; email=; avatar=; name=; expires=Thu, 01 Jan 1970 00:00:01 GMT ;";
    $("#customBtn")[0].style.display = "inline-block";
  });
}
