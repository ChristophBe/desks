<template>
  <v-table>
    <thead>
    <tr>
      <th class="text-left w-33">
        Room
      </th>
      <th class="text-left">
        Date
      </th>
      <th class="text-left">
        Start
      </th>
      <th class="text-left">
        End
      </th>
      <th></th>
    </tr>
    </thead>
    <tbody>

    <tr
        v-for="item in bookings"
        :key="item.id"
    >
      <td>{{ item.room.name }}</td>
      <td>{{ formatDate(item.start) }}</td>
      <td>{{ formatTime(item.start) }}</td>
      <td>{{ formatTime(item.end) }}</td>
      <td class="d-flex justify-end align-center">
        <v-btn
            icon
            elevation="0"
        >
          <v-icon>mdi-dots-vertical</v-icon>

          <v-menu activator="parent">
            <v-list density="compact">
              <v-list-item class="align-center" @click="editBooking(item)" :disabled="isCurrentOrPastBooking(item)">
                <v-list-item-avatar start icon="mdi-pencil"></v-list-item-avatar>
                <v-list-item-title>Edit</v-list-item-title>
              </v-list-item>
              <delete-booking-menu-item :booking="item" :disabled="isCurrentOrPastBooking(item)"></delete-booking-menu-item>
            </v-list>
          </v-menu>
        </v-btn>
      </td>
    </tr>
    </tbody>
  </v-table>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import moment from "moment";
import DeleteBookingMenuItem from "@/components/booking-components/DeleteBookingMenuItem.vue";
import Booking from "@/models/Booking";

export default defineComponent({
  name: "BookingsTable",
  components: {DeleteBookingMenuItem},
  // type inference enabled
  props: {
    bookings: Array
  },
  methods: {
    formatDate(date:Date) {
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date:Date) {
      return moment(date).format("HH:mm")
    },

    editBooking(booking:Booking){
      this.$emit("editBooking", booking)
    },
    isCurrentOrPastBooking(booking:Booking):boolean{
      return moment(booking?.start).isBefore(moment.now())
    },
  }

})
</script>

<style scoped>

</style>
