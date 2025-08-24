<template>
  <div class="event-list">
    <h2>Events</h2>
    <div class="event-form">
      <input v-model="newEvent.name" placeholder="Event name" />
      <input v-model="newEvent.description" placeholder="Description" />
      <input v-model="newEvent.startTime" type="datetime-local" />
      <input v-model="newEvent.endTime" type="datetime-local" />
      <button @click="createEvent">Create Event</button>
    </div>
    <ul class="event-items">
      <li v-for="event in events" :key="event.id" class="event-item">
        <div class="event-info">
          <h3>{{ event.name }}</h3>
          <p>{{ event.description }}</p>
          <p>Start: {{ formatDate(event.startTime) }}</p>
          <p>End: {{ formatDate(event.endTime) }}</p>
        </div>
        <div class="event-actions">
          <button @click="editEvent(event)">Edit</button>
          <button @click="deleteEvent(event.id)">Delete</button>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Event {
  id: string
  name: string
  description: string
  startTime: string
  endTime: string
}

const events = ref<Event[]>([])
const newEvent = ref({
  name: '',
  description: '',
  startTime: '',
  endTime: ''
})

const API_BASE = 'http://localhost:8080/api'

const fetchEvents = async () => {
  try {
    const response = await fetch(`${API_BASE}/events`)
    if (response.ok) {
      events.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to fetch events:', error)
  }
}

const createEvent = async () => {
  try {
    const response = await fetch(`${API_BASE}/events`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newEvent.value)
    })
    if (response.ok) {
      newEvent.value = { name: '', description: '', startTime: '', endTime: '' }
      fetchEvents()
    }
  } catch (error) {
    console.error('Failed to create event:', error)
  }
}

const editEvent = (event: Event) => {
  // For simplicity, just log the event to edit
  console.log('Edit event:', event)
  // In a real implementation, you'd open a modal or form for editing
}

const deleteEvent = async (id: string) => {
  try {
    const response = await fetch(`${API_BASE}/events/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      fetchEvents()
    }
  } catch (error) {
    console.error('Failed to delete event:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  fetchEvents()
})
</script>

<style scoped>
.event-list {
  margin: 20px;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

.event-form {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.event-form input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.event-form button {
  padding: 8px 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.event-items {
  list-style: none;
  padding: 0;
}

.event-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.event-info h3 {
  margin: 0 0 5px 0;
}

.event-actions {
  display: flex;
  gap: 10px;
}

.event-actions button {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.event-actions button:first-child {
  background-color: #2196F3;
  color: white;
}

.event-actions button:last-child {
  background-color: #f44336;
  color: white;
}
</style>