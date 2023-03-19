<template>
  <v-menu :activator="activator">
    <v-list density="compact">
      <v-list-item class="align-center" @click="onEdite()" :disabled="isAllowedToEdit()">
        <template v-slot:prepend>
          <v-icon>mdi-pencil</v-icon>
        </template>
        <v-list-item-title>Edit</v-list-item-title>
      </v-list-item>
      <v-list-item class="align-center" @click="onDelete()" :disabled="isCurrentOrPastBooking()">
        <template v-slot:prepend>
          <v-icon>mdi-delete</v-icon>
        </template>
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
    },
    isAllowedToEdit(): boolean {
      return moment(this.booking?.end).isBefore(moment.now())
    }
  }
});
</script>

<style scoped>

</style>
