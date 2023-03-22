<template>
<v-row>
  <v-col
      v-for="item in [...bookings].sort(compareBookingsByTime)"
      :key="item.id"
    :style="{maxWidth: '478px'}">
    <v-card
        :variant="'flat'"
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
        {{ dateConverter(item.start).format("dddd") }}
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
  <!-- <v-table>
    <thead>
    <tr>
      <th class="text-left w-33">
        Room
      </th>
      <th class="text-left w-25">
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
        <v-btn v-if="isOngoing(item)" @click.stop.prevent="() => leaveEarly(item)" variant="flat">Leave now</v-btn>
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
-->
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import BookingContextMenu from "@/components/booking-components/BookingContextMenu.vue";
import moment from "moment/moment";
import {mapActions} from "vuex";
import {Moment} from "moment";

export default defineComponent({
  name: "BookingsTable",
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
    async onDelete() {
      await this.deleteBooking(this.booking);
      this.$emit("deleted")
    },

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
    },
    isCurrentOrPastBooking(): boolean {
      return moment(this.booking?.start).isBefore(moment.now())
    },
    isAllowedToEdit(): boolean {
      return moment(this.booking?.end).isBefore(moment.now())
    }
  }

})
</script>

<style scoped>

</style>
