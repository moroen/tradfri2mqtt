import { Notify } from "quasar";
import { newWebSocket, sendWSCommand } from "./websocket";

export const showError = function (msg) {
  Notify.create(
    {
      type: "negative",
      message: msg,
      position: "center",
      progress: true,

      actions: [
        {
          label: "Dismiss",
          color: "white",
          handler: () => {
            /* ... */
          },
        },
      ],
    },
    2000
  );
};

const startWebSocket = (context) => {
  console.log("startWebSocket");
  const connection = newWebSocket(context);
  context.commit("setConnection", connection);
};

const getSettings = (context) => {
  console.log("getSettings");
  var settingsURL = "api/settings";

  fetch(settingsURL)
    .then((response) => {
      // console.log(response)
      return response.json();
    })

    .then((data) => {
      // console.log(data)
      context.commit("setConfig", data);
    })
    .catch((err) => {
      console.error("Unable to get settings - connection failed");
      showError("Unable to get settings");
      // console.error(err);
    });
};

const getNewPSK = (context, payload) => {
  var settingsURL = "api/getPSK";

  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  };

  // if (import.meta.env.MODE == "development") {
  //  settingsURL = "http://localhost:8321/api/getPSK";
  //}

  fetch(settingsURL, requestOptions)
    .then((response) => {
      if (response.status == 200) {
        return response.json();
      } else {
        throw new Error("Unable");
      }
    })
    .then((data) => {
      context.commit("setPSK", {
        identity: data.identity,
        passkey: data.passkey,
      });
    })
    .catch((err) => {
      showError("Unable to generate new PSK-pair");
    });
};

const saveSettings = (context) => {
  var settingsURL = "api/settings";

  console.log("saveSettings action");

  const settings = context.getters.SettingsJson;

  // if (import.meta.env.MODE == "development") {
  //  settingsURL = "http://localhost:8321/api/settings";
  //}

  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: settings,
  };

  fetch(settingsURL, requestOptions)
    .then((response) => {
      console.log("Got response");
      if (response.status == 200) {
        return response.json();
      } else {
        throw new Error("Unable");
      }
    })
    .then((data) => {
      console.log("Got data");
      console.log(data);
    })
    .catch((err) => {
      console.error("Unable to save setting\n" + err);
      showError("Unable to save settings");
    });
};

const doTest = () => {
  console.log("Do test");
  showError("Testing");
};

export default {
  getSettings,
  doTest,
  startWebSocket,
  sendWSCommand,
  getNewPSK,
  saveSettings,
};
