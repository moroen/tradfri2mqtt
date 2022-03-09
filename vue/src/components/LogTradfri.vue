<template>
  <div>
    <q-checkbox v-model="showLog" @update:model-value="onLogToggle">Show log</q-checkbox>
  </div>
  <scroller
    ref="scrollArea"
    :expand="true"
    :bottom-margin="20"
    height="200px"
    style="width: 99%"
  >
    <div v-for="(item, index) in items" :key="index" class="row q-gutter-md">
      <div class="col-auto">{{ item.time }}</div>
      <div class="col-auto" :class="levelColor(item.level)">{{ item.level }}</div>
      <div class="col-auto">{{ item.message }}</div>
    </div>
  </scroller>
</template>

<script setup>
import { useStore } from "vuex";
import { computed, watch, ref } from "vue";
import { mapFields } from "vuex-map-fields";
import scroller from "./Scroller.vue";

const store = useStore();

const scrollArea = ref(null);

const levelColor = (level) => {
  if (level == "error") {
    return "text-negative";
  } else {
    return "text-positive";
  }
};

const onLogToggle = function (value) {
  console.log("Hi");
  var cmd;
  if (value) {
    cmd = "start";
  } else {
    cmd = "stop";
  }
  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "log",
      command: cmd,
    })
  );
  if (!value) {
    store.commit("clearLog");
  }
};

const items = computed({
  get() {
    return store.state.websocket.log;
  },
});

watch(
  items,
  (oldValue, newValue) => {
    console.log("Changed");
    scrollArea.value.scrollToBottom();
  },
  { deep: true }
);
</script>

<script>
export default {
  computed: {
    ...mapFields(["websocket.log", "websocket.showLog"]),
  },
};
</script>
