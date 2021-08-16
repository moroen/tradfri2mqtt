<template>
  <q-list>
    <q-item to="/tradfri">
      <q-item-section avatar>
        <q-icon name="lightbulb" />
      </q-item-section>

      <q-item-section>
        <q-item-label>Tradfri</q-item-label>
        <q-item-label caption>
          Manage IKEA-Gateway and devices
        </q-item-label>
      </q-item-section>
    </q-item>

    <q-item to="/mqtt">
      <q-item-section avatar>
        <q-icon name="share" />
      </q-item-section>

      <q-item-section>
        <q-item-label>MQTT</q-item-label>
        <q-item-label caption>
          MQTT settings and tools
        </q-item-label>
      </q-item-section>
    </q-item>

    <q-item to="/interface">
      <q-item-section avatar>
        <q-icon name="settings" />
      </q-item-section>

      <q-item-section>
        <q-item-label>Interface</q-item-label>
        <q-item-label caption>
          Interface settings
        </q-item-label>
      </q-item-section>
    </q-item>

    <q-item to="/about">
      <q-item-section avatar>
        <q-icon name="info" />
      </q-item-section>

      <q-item-section>
        <q-item-label>About</q-item-label>
        <q-item-label caption>
        </q-item-label>
      </q-item-section>
    </q-item>

  </q-list>
</template>

<script>
export default {
  name: "Navigation",
  setup() {
    return {};
  }
};
</script>
