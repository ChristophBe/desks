<template>

  <v-list-item class="align-center" @click="onClick" :disabled="disabled">

    <v-list-item-avatar start icon="mdi-delete"></v-list-item-avatar>
    <v-list-item-title>Remove</v-list-item-title>
  </v-list-item>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {mapActions} from "vuex";

export default defineComponent({
  name: "delete-booking-menu-item",
  props:{
    booking:{
      type: Object
    },
    disabled: {
      type: Boolean
    }
  },
  computed: {

  },
  methods: {
    ...mapActions("bookings", ["deleteBooking"]),
    async onClick() {
      await this.deleteBooking(this.booking);
    }
  }
});
</script>

<style scoped>

</style>
