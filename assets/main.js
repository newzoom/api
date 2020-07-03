let u;
let displaySignIn = true;
function getUserLocalData() {
  let retrievedData = localStorage.getItem("user");
  if (retrievedData != undefined && retrievedData != null) {
    u = JSON.parse(retrievedData);
    displaySignIn = false;
  }
}
getUserLocalData();
