<template>
<q-form>
          <q-checkbox v-model="enable" label="Enable Tradfri" />
          <q-input
            filled
            v-model="gateway"
            label="Gateway IP *"
            hint="Hostname or IP-address of the IKEA-tradfri gateway"
            :rules="stringRules"
            :disable="!enable"
          />

          <div class="row items-center content-stretch">
            <div class="col-8">
              <q-input
                filled
                v-model="identity"
                label="Ident *"
                hint="Ident"
                :rules="stringRules"
                :disable="!enable"
              />
              <q-input
                filled
                :disable="!enable"
                v-model="passkey"
                label="Key *"
                hint="PSK"
                :rules="stringRules"
              />
            </div>
            <div class="q-pa-md">
              <q-btn
                label="Generate ID/PSK"
                color="primary"
                @click="promptPSK"
                :disable="!enable"
              ></q-btn>
            </div>
          </div>
          <div class="row items-center">
            <div class="col-3">
              <q-input
                filled
                :disable="!enable"
                v-model.number="keepalive"
                label="Keep Alive interval:"
                hint="Dis- and reconnect at every n seconds. 0 to disable"
                :rules="timerRule"
                type="number"
              />
            </div>
            <div class="col-2"><q-space></q-space></div>
            <div class="col-3">
              <q-input
                filled
                :disable="!enable"
                v-model.number="disconnecttimer"
                label="Disconnect timer:"
                hint="Disconnect after n seconds. 0 to disable"
                :rules="timerRule"
                type="number"
              />
            </div>
          </div>
        </q-form>
</template>

<script setup>
import { useQuasar } from "quasar";
import { useStore } from "vuex";
import { ref } from "vue";

import { stringRules, timerRule } from "../rules.js";

import getPSK from "../components/GetPSK.vue";
import { mapFields } from "vuex-map-fields";

const store = useStore();
const q = useQuasar();

const tab = ref("devices");

// const nameRules = [val => (val && val.length > 0) || 'Please type something']

function promptPSK() {
  q.dialog({
    component: getPSK,
  })
    .onOk(() => {
      // console.log("OK");
    })
    .onCancel(() => {
      // console.log("Cancel");
    });
}
</script>

<script>
export default {
  computed: {
    ...mapFields([
      "settings.tradfri.gateway",
      "settings.tradfri.identity",
      "settings.tradfri.passkey",
      "settings.tradfri.enable",
      "settings.tradfri.keepalive",
      "settings.tradfri.disconnecttimer",
    ]),
  },
};
</script>
