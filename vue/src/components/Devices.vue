<template>
  <div ref="refToolbar">
    <q-toolbar>
      <q-space />
      <q-btn label="Test" color="primary" @click="test" ref="input"></q-btn>
      <q-btn label="Update" color="primary" @click="UpdateDevices"></q-btn>
    </q-toolbar>
  </div>

  <div ref="refScroller">
    <q-virtual-scroll :visible="true" :items="devices" class="Scroll">
      <template v-slot="{ item, index }">
        <q-item :key="index" dense>
          <q-item-section>
            <q-item-label>{{ item }}</q-item-label>
          </q-item-section>
        </q-item>
      </template>
    </q-virtual-scroll>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { mapFields } from "vuex-map-fields";
import { useStore } from "vuex";

const store = useStore();

const refToolbar = ref(null);
const refScroller = ref(null);

const scrollColor = ref("red");
const scrollHeight = ref("100px");

let scrollTopPosition = 0;

const resizeScroll = () => {
  scrollColor.value = "green";
  scrollHeight.value = heighFn();
};

const heighFn = () => {
  let h = window.innerHeight - scrollTopPosition - 20;
  return h + "px";
};

onMounted(() => {
  scrollTopPosition = refScroller.value.getBoundingClientRect().top;
  resizeScroll();
  window.addEventListener("resize", resizeScroll);
});

const test = () => {};

const UpdateDevices = () => {
  console.log("UpdateDevices");
  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "devices",
      command: "update",
    })
  );
};
</script>

<script>
export default {
  computed: {
    ...mapFields(["websocket.devices"]),
  },
  mounted() {},
};
</script>

<style>
.Scroll {
  background-color: bisque;
  color: v-bind("scrollColor");
  height: v-bind("scrollHeight");
  border: 1px solid black;
}
</style>
