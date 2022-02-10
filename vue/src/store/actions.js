import { Notify } from "quasar";

const showError = function(msg) {
    Notify.create({
      type: 'negative',
      message: msg,
      progress: true,
  
      actions: [
        { label: 'Dismiss', color: 'white', handler: () => { /* ... */ } }
      ]
    }, 2000)
  }

const getSettings = (context) => {
  var settingsURL = "api/settings";

  if (import.meta.env.MODE == "development") {
    settingsURL = "http://localhost:8321/api/settings";
  }
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
      context.commit("setStatus", "Unable to get settings");
      showError("Unable to get settings");
      console.error(err);
    });
};

const getNewPSK = (context, payload) => {
  var settingsURL = "api/getPSK";

  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  };

  if (import.meta.env.MODE == "development") {
    settingsURL = "http://localhost:8321/api/getPSK";
  }

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

  const settings = context.getters.SettingsJson;

  if (import.meta.env.MODE == "development") {
    settingsURL = "http://localhost:8321/api/settings";
  }

  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: settings,
  };

  fetch(settingsURL, requestOptions)
  .then((response) => {
    if (response.status == 200) {
      return response.json();
    } else {
      throw new Error("Unable");
    }
  })
  .then((data) => {
    console.log(data)
  })
  .catch((err) => {
    showError("Unable to save settings");
  });

  console.log(settings);
};

export default {
  getSettings,
  getNewPSK,
  saveSettings,
};
