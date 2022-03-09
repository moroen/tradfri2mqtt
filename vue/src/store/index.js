import { Notify } from "quasar";
import { createStore } from "vuex";
import { ref } from "vue";
import actions from "./actions";

import { getField, updateField } from "vuex-map-fields";
import {
  addLogEntry,
  addWSLogEntry,
  updateDeviceInfo,
  clearLog,
} from "./websocket";

export default createStore({
  state: {
    settings: {
      tradfri: {
        enable: true,
        gateway: "127.0.0.2",
        identity: "abcde",
        passkey: "12345",
        keepalive: ref(0),
        disconnecttimer: ref(0),
      },
      mqtt: {
        commandtopic: "tradfri",
        discoverytopic: "homeassistant",
        enable: true,
        host: "127.0.0.1",
        port: 1883,
      },

      interface: {
        enable: true,
        port: 8321,
        serverroot: "./www",
      },
    },
    footBar: ref(null),
    status: "Ok",
    websocket: {
      connection: ref(null),
      log: [],
      wslog: [],
      showLog: false,
      devices: [],
    },
  },
  mutations: {
    updateField,
    addLogEntry,
    addWSLogEntry,
    updateDeviceInfo,
    clearLog,
    setFootBarRef(state, payload) {
      state.footBar = payload;
    },
    setConnection(state, payload) {
      state.websocket.connection = payload;
    },
    setConfig(state, payload) {
      state["settings"][payload.section][payload.key] = payload.value;
    },
    setTradfri(state, payload) {
      // console.log(payload)
      state.settings.tradfri = payload;
    },
    setConfig(state, payload) {
      // console.log(payload)
      state.settings.tradfri = payload.tradfri;
      state.settings.mqtt = payload.mqtt;
      state.settings.interface = payload.interface;
    },
    setPSK(state, payload) {
      // console.log(payload)
      state.settings.tradfri.identity = payload.identity;
      state.settings.tradfri.passkey = payload.passkey;
    },

    setStatus(state, status) {
      state.status = status;
    },
  },
  actions,
  modules: {},
  getters: {
    getField,
    SettingsJson: (state) => {
      return JSON.stringify(state.settings);
    },
    wsConnection: (state) => {
      return state.websocket.connection;
    },
  },
});
