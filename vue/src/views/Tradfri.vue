<template>
  <q-page padding>
    <q-tabs
      v-model="tab"
      dense
      class="text-grey"
      active-color="primary"
      indicator-color="primary"
      align="left"
      narrow-indicator
    >
      <q-tab name="gateway" label="Gateway" />
      <q-tab name="devices" label="Devices" />
    </q-tabs>

    <q-separator />

    <q-tab-panels v-model="tab" animated>
      <q-tab-panel name="gateway">
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
      </q-tab-panel>

      <q-tab-panel name="devices">
        <div class="text-h6">Devices</div>
        Lorem ipsum dolor sit amet consectetur adipisicing elit.
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup>
import { useQuasar } from "quasar";
import { useStore } from "vuex";
import { computed, ref } from "vue";

import { stringRules, timerRule } from "../rules.js";

import getPSK from "../components/GetPSK.vue";
import { mapFields } from "vuex-map-fields";

const store = useStore();
const q = useQuasar();

const tab = ref("gateway");

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
      "tradfri.gateway",
      "tradfri.identity",
      "tradfri.passkey",
      "tradfri.enable",
      "tradfri.keepalive",
      "tradfri.disconnecttimer",
    ]),
  },
};
</script>
