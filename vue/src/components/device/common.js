const sendDeviceCommandObject = (deviceID, command) => {
  const uri = `/api/devices/${deviceID}/set`;

  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(command),
  };

  fetch(uri, requestOptions)
    .then((response) => {
      if (response.status == 200) {
        return response.json();
      } else {
        throw new Error("Unable to set state");
      }
    })
    .then((data) => {
      // console.log("Got data");
      // console.log(data);
    })
    .catch((err) => {
      console.error("Unable to save setting\n" + err);
      showError("Unable to save settings");
    });
};

export default sendDeviceCommandObject;

const sendDeviceWSCommand = (deviceID, command);
