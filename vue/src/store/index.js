import { Notify } from 'quasar'
import { createStore } from 'vuex'
import { ref } from 'vue'
import actions from './actions';

import { getField, updateField } from "vuex-map-fields";



export default createStore({
  state: {
    tradfri: {
      enable: true,
      gateway: "127.0.0.2",
      identity: "abcde",
      passkey: "12345",
      keepalive: ref(0),
      disconnecttimer: ref(0)
    },
    mqtt: {
      commandtopic:	"tradfri",
      discoverytopic:	"homeassistant",
      enable:	true,
      host:	"127.0.0.1",
      port: 1883
    },

    interface: {
      enable: true,
      port :8321,
      serverroot: "./www"
    },
    status: "Ok"
  },
  mutations: {
    updateField,
    setConfig (state, payload) {
      state[payload.section][payload.key] = payload.value
    },
    setTradfri (state, payload) {
      // console.log(payload)
      state.tradfri = payload
    },
    setConfig (state, payload) {
      // console.log(payload)
      state.tradfri = payload.tradfri
      state.mqtt = payload.mqtt
      state.interface = payload.interface
    },
    setPSK(state, payload) {
      // console.log(payload)
      state.tradfri.identity = payload.identity
      state.tradfri.passkey = payload.passkey
    },

    setStatus(state, status) {
      state.status=status
    }
  },
  actions,
  modules: {
  },
  getters: {
    getField,
    SettingsJson: (state) => {
      return JSON.stringify(state)
  }
  }
})
