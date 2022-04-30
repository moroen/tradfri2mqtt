<template>
    <div class="column">
        <div class="col self-end">{{ level }}</div>
        <div>
            <q-slider v-model="level" :min="0" :max="100" @change="saveValue" />
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useStore } from "vuex";

const store = useStore();

const props = defineProps({
    deviceid: {
        required: true,
    },
});

const saveValue = () => {
    const command = {
        class: "device",
        // command: "state",
        value: {
            deviceid: props.deviceid,
            payload: { position: store.state.websocket.devices[props.deviceid].position },
        }
    };
    store.dispatch("sendWSCommand", JSON.stringify(command));
};

const level = computed({
    get() {
        return store.state.websocket.devices[props.deviceid].position;
    },
    set(value) {
        var payload = {
            deviceid: props.deviceid,
            level: value,
        };
        store.commit("setDeviceBlindPosition", payload);
    },
});

const standard = ref(0);
</script>
