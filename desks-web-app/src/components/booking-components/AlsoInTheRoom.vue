<template>
  <div>
    <div v-if="bookingsOfColleagues === undefined || isLoadingBookingsByRoomAndDay(roomId, date)">
      <v-progress-circular class="mx-4 mr-4" indeterminate size="small"></v-progress-circular>
      load bookings of your colleagues
    </div>
    <v-alert v-else-if="bookingsOfColleagues.length <= 0" class="mx-4">There are no other bookings for the same
      office.
    </v-alert>

    <v-list v-else>
      <v-list-item two-line v-for="booking in bookingsOfColleagues.sort(compareByName)" :key="booking.id">
        <v-list-item-avatar start>
          <v-avatar color="primary">
            {{ booking.user.givenName.charAt(0) }}{{ booking.user.familyName.charAt(0) }}
          </v-avatar>


        </v-list-item-avatar>
        <v-list-item-header>
          <v-list-item-title>{{ booking.user.givenName }} {{ booking.user.familyName }}</v-list-item-title>
          <v-list-item-subtitle>{{ $format.timeRange(booking.start, booking.end) }}</v-list-item-subtitle>
        </v-list-item-header>

      </v-list-item>
    </v-list>
  </div>

</template>

<script lang="ts">

import {defineComponent} from "vue";
import Booking from "@/models/Booking";
import {mapActions, mapGetters, mapState} from "vuex";
import moment from "moment";

interface alsoInTheRoomData {
  bookingsOfColleagues?: Array<Booking>
}

export default defineComponent({
  name: "AlsoInTheRoom",
  data: (): alsoInTheRoomData => ({
    bookingsOfColleagues: undefined,
  }),
  // type inference enabled
  props: {
    roomId: {
      type: Number
    },
    density: String,
    date: String
  },

  async mounted() {
    this.fetchData()
    this.updateBookingsOfColleagues();
  },
  computed: {
    ...mapState('user', ['user']),
    ...mapState('bookings', ['bookings']),
    ...mapGetters('bookings', ['getBookingsByRoomAndDay', 'isLoadingBookingsByRoomAndDay'])
  },
  watch: {
    roomId(){
      this.fetchData()
    },
    date(){
      this.fetchData()
    },
    bookings() {
      this.updateBookingsOfColleagues();
    }
  },
  methods: {
    ...mapActions('bookings', ['fetchBookingsByRoomAndDate']),
    updateBookingsOfColleagues() {
      if (!this.user) {
        return
      }
      this.bookingsOfColleagues = this.getBookingsByRoomAndDay(this.roomId, this.date)
          .filter((booking: Booking) => booking.user.id !== this.user.id)

    },
    compareByName(a: Booking, b: Booking): number {
      let comparison = a.user.familyName.localeCompare(b.user.familyName)
      if (comparison === 0) {
        comparison = a.user.givenName.localeCompare(b.user.givenName)
      }
      if (comparison === 0) {
        comparison = moment(a.start).diff(moment(b.start))
      }
      return comparison
    },
    fetchData(){
      this.bookingsOfColleagues = undefined
      if(this.roomId && this.date){
        this.fetchBookingsByRoomAndDate({roomId: this.roomId, date: this.date})
      }
      this.updateBookingsOfColleagues();
    }
  }
})
</script>

<style scoped>

</style>
