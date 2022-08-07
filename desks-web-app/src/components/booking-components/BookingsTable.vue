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
        v-for="item in [...bookings].sort(compareBookingsByTime)"
        :key="item.id"
        @click="$emit('openBooking', item)"
    >
      <td>{{ item.room.name }}</td>
      <td>{{ $format.date(item.start) }}</td>
      <td>{{ $format.time(item.start) }}</td>
      <td>{{ $format.time(item.end) }}</td>
      <td class="d-flex justify-end align-center">
        <v-btn
            icon
            elevation="0"
        >
          <v-icon>mdi-dots-vertical</v-icon>
          <booking-context-menu activator="parent" :booking="item" @edit="editBooking(item)"></booking-context-menu>
        </v-btn>
      </td>
    </tr>
    </tbody>
  </v-table>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import BookingContextMenu from "@/components/booking-components/BookingContextMenu.vue";
import moment from "moment/moment";

export default defineComponent({
  name: "BookingsTable",
  components: {BookingContextMenu},
  // type inference enabled
  props: {
    bookings: Array
  },
  methods: {
    editBooking(booking: Booking) {
      this.$emit("editBooking", booking)
    },

    compareBookingsByTime(a: Booking, b: Booking): number {
      return moment(a.start).diff(moment(b.start))
    }
  }

})
</script>

<style scoped>

</style>
