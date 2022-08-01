<template>
  <v-menu :activator="activator">
    <v-list density="compact">
      <v-list-item class="align-center" @click="onEdite()" :disabled="isCurrentOrPastBooking()">
        <v-list-item-avatar start icon="mdi-pencil"></v-list-item-avatar>
        <v-list-item-title>Edit</v-list-item-title>
      </v-list-item>
      <v-list-item class="align-center" @click="onDelete()" :disabled="isCurrentOrPastBooking()">

        <v-list-item-avatar start icon="mdi-delete"></v-list-item-avatar>
        <v-list-item-title>Remove</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import {mapActions} from "vuex";
import moment from "moment/moment";

export default defineComponent({
  name: "booking-context-menu",
  props: {
    booking: {
      type: Object
    },
    activator: String
  },

  methods: {
    ...mapActions("bookings", ["deleteBooking"]),
    async onDelete() {
      await this.deleteBooking(this.booking);
      this.$emit("deleted")
    },
    onEdite() {
       this.$emit("edit")
    },
    isCurrentOrPastBooking(): boolean {
      return moment(this.booking?.start).isBefore(moment.now())
    }
  }
});
</script>

<style scoped>

</style>
