<template>
  <!-- notice dialogRef here -->
  <q-dialog ref="dialogRef" @hide="onDialogHide">
    <q-card class="q-dialog-plugin">
      <q-input
        ref="gateway"
        filled
        v-model="gateway"
        label="Gateway IP *"
        hint="Hostname or IP-address of the IKEA-tradfri gateway"
        lazy-rules
        :rules="[(val) => (val && val.length > 0) || 'Please type something']"
      />
      <q-input
        filled
        v-model="masterKey"
        label="Master key *"
        hint="Key found under the Tradfri Gateway"
        lazy-rules
        :rules="[(val) => (val && val.length > 0) || 'Please type something']"
      />

      <q-card-actions align="right">
        <q-btn color="primary" label="OK" @click="onOKClick" />
        <q-btn color="primary" label="Cancel" @click="onCancelClick" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script>
import { useDialogPluginComponent } from "quasar";
import { ref } from "vue";
import { computed } from "vue";
import { useStore } from "vuex";

import { mapFields } from "vuex-map-fields";
import store from "../store";
import TradfriVue from "../views/Tradfri.vue";

export default {
  methods: {
    show() {
      this.$refs.dialogRef.show();
      return this.dialogRef;
    },
  },

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

  props: {
    // ...your custom props
  },

  emits: [
    // REQUIRED; need to specify some events that your
    // component will emit through useDialogPluginComponent()
    ...useDialogPluginComponent.emits,
  ],

  setup() {
    const $store = useStore();
    const masterKey = ref(null);

    // REQUIRED; must be called inside of setup()
    const {
      dialogRef,
      onDialogHide,
      onDialogOK,
      onDialogCancel,
    } = useDialogPluginComponent();
    // dialogRef      - Vue ref to be applied to QDialog
    // onDialogHide   - Function to be used as handler for @hide on QDialog
    // onDialogOK     - Function to call to settle dialog with "ok" outcome
    //                    example: onDialogOK() - no payload
    //                    example: onDialogOK({ /*.../* }) - with payload
    // onDialogCancel - Function to call to settle dialog with "cancel" outcome

    return {
      masterKey,
      // This is REQUIRED;
      // Need to inject these (from useDialogPluginComponent() call)
      // into the vue scope for the vue html template
      dialogRef,
      onDialogHide,

      // other methods that we used in our vue html template;
      // these are part of our example (so not required)
      onOKClick() {

        const postData = {
          gateway: $store.state.tradfri.gateway,
          key: masterKey.value,
        };

        // console.log(postData);
        $store.dispatch("getNewPSK", postData);
        // on OK, it is REQUIRED to
        // call onDialogOK (with optional payload)
        onDialogOK();
        // or with payload: onDialogOK({ ... })
        // ...and it will also hide the dialog automatically
      },

      // we can passthrough onDialogCancel directly
      onCancelClick: onDialogCancel,
    };
  },
};
</script>
