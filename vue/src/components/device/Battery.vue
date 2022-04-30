<template>
    <div class="items-start" v-if="hasBattery()">
        <q-icon name="battery_3_bar" />
        {{ status }}
    </div>
</template>

<script setup>
import { computed } from "vue";
import { useStore } from "vuex";

const store = useStore();

const props = defineProps({
    deviceid: {
        required: true,
    },
});

const hasBattery = () => {
    return (store.state.websocket.devices[props.deviceid].battery != -1)
}

const status = computed({
    get() {
        return store.state.websocket.devices[props.deviceid].battery;
    },
});
</script>