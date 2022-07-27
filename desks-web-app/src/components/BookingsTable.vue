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
              <delete-booking-menu-item :booking="item"></delete-booking-menu-item>
            </v-list>
          </v-menu>
        </v-btn>
      </td>
    </tr>
    </tbody>
  </v-table>
</template>

<script type="ts">
import {defineComponent} from "vue";
import moment from "moment";
import DeleteBookingMenuItem from "@/components/DeleteBookingMenuItem.vue";

export default defineComponent({
  name: "BookingsTable",
  components: {DeleteBookingMenuItem},
  // type inference enabled
  props: {
    bookings: Array
  },

  methods: {
    formatDate(date) {
      return moment(date).format("DD.MM.YYYY")
    },
    formatTime(date) {
      return moment(date).format("HH:mm")
    },
  }

})
</script>

<style scoped>

</style>
