<template>
<v-row class="bookingsRow justify-start flex-0-0">
  <v-col
      v-for="item in [...bookings].sort(compareBookingsByTime)"
      :key="item.id"
    class="bookingsCol">
    <v-card
        :variant="'flat'"
        class="card"
    >
      <template v-slot:append>
        <v-btn icon elevation="0" >
          <v-icon>mdi-dots-vertical</v-icon>
          <booking-context-menu activator="parent" :booking="item" @edit="editBooking(item)"></booking-context-menu>
        </v-btn>
      </template>
      <template v-slot:title v-if="dateConverter(item.start).startOf('day').isSame(moment().startOf('day'))">
        Today
      </template>
      <template v-slot:title v-else-if="dateConverter(item.start).isSameOrBefore(dateConverter(item.start).add(1, 'week'))">
        {{ dateConverter(item.start).format("DD.MM.YYYY") }}
      </template>
      <template v-slot:title v-else>
        {{ $format.date(item.start) }}
      </template>
      <template v-slot:subtitle>
        {{ item.room.name }}
      </template>


      <v-card-text>

         <v-chip>{{ $format.time(item.start) }} - {{ $format.time(item.end) }}</v-chip>

      </v-card-text>
    </v-card>
  </v-col>
</v-row>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import BookingContextMenu from "@/components/booking-components/BookingContextMenu.vue";
import moment from "moment/moment";
import {mapActions} from "vuex";
import {Moment} from "moment";

export default defineComponent({
  name: "BookingsCards",
  components: {BookingContextMenu},
  // type inference enabled
  props: {
    bookings: Array
  },
  methods: {
    ...mapActions("bookings", ["leaveEarly", "deleteBooking"]),
    editBooking(booking: Booking) {
      this.$emit("editBooking", booking)
    },
    moment: () => moment(),

    compareBookingsByTime(a: Booking, b: Booking): number {
      return moment(a.start).diff(moment(b.start))
    },

    isOngoing(booking:Booking): boolean{
      return moment().isBetween(booking.start,booking.end)
    },
    dateConverter(d: Date): Moment {
      return moment(d);
    },
    onEdite() {
      this.$emit("edit")
    }
  }

})
</script>

<style>

@media screen and (min-device-width: 970px) {
  .bookingsRow {
    min-width: 20%;
    width: fit-content;
  }
}

.bookingsCol {
  flex-grow: 0;
}

.card {
  width: max-content;
}
</style>
