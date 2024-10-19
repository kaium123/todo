<template>
  <div class="todo-main">
    <h1>TODOリスト</h1>

    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>

    <div class="input-group">
      <input
        v-model="newTask"
        placeholder="新しいタスクを入力 (Enter new task)"
        @keyup.enter="addTodo"
      />
      <select v-model="newPriority" class="priority-dropdown">
        <option value="high">高 (High)</option>
        <option value="medium">中 (Medium)</option>
        <option value="low">低 (Low)</option>
      </select>
      <button class="green-button" @click="addTodo">追加 (Add)</button>
    </div>

    <div v-if="filteredTodos.length > 0">
      <div v-for="todo in filteredTodos" :key="todo.ID" class="todo-item">
        <div>
          <span :class="{ 'done-task': todo.Status === 'done' }">
            {{ todo.Task }} (優先度: {{ todo.Priority }})
          </span>
        </div>
        <div class="buttons">
          <button class="edit-button" @click="enableEdit(todo)" v-if="!todo.isEditing">編集 (Edit)</button>
          <button class="delete-button" @click="deleteTodo(todo.ID)" v-if="!todo.isEditing">🗑️</button>
        </div>

        <!-- Edit dropdown section -->
        <div v-if="todo.isEditing" class="edit-dropdown">
          <div class="dropdown-container">
            <!-- Task input aligned to the left -->
            <input
              v-if="todo.editOption === 'task'"
              v-model="todo.Task"
              class="edit-input"
              @blur="editTodo(todo)"
              @keyup.enter="editTodo(todo)"
              placeholder="Edit Task"
            />

            <select v-model="todo.editOption" @change="updateEditOption(todo)" class="right-align">
              <option value="task">タスク (Task)</option>
              <option value="status">ステータス (Status)</option>
            </select>

            <select
              v-if="todo.editOption === 'status'"
              v-model="todo.Status"
              class="status-dropdown"
              @blur="editTodo(todo)"
              @keyup.enter="editTodo(todo)"
            >
              <option value="created">作成済み (Created)</option>
              <option value="processing">処理中 (Processing)</option>
              <option value="done">完了 (Done)</option>
            </select>
          </div>
        </div>
      </div>
    </div>
    <div v-else>
      <p>タスクがありません (No tasks available).</p>
    </div>
  </div>
</template>


<script>
export default {
  data() {
    return {
      newTask: '',
      newPriority: 'low',
      todos: [],
      statusMessage: '',
      searchType: 'task',
      searchTask: '',
      searchStatus: '',
      filteredTodos: [], // Introduced filteredTodos for search results
    };
  },
  mounted() {
    this.fetchTodos();
  },
  methods: {
    async fetchTodos() {
      try {
        const response = await fetch(`/api/v1/todos`);
        if (!response.ok) throw new Error(`Failed to get todo list. statusCode: ${response.status}`);
        
        const data = await response.json();
        this.todos = data.data;
        this.filteredTodos = this.todos; // Initialize filteredTodos
        // Remove or comment out this line to prevent displaying the success message
        // this.statusMessage = 'タスクが取得されました (Tasks fetched successfully)'; 
      } catch (error) {
        console.error(error);
        this.statusMessage = 'タスクの取得に失敗しました (Failed to fetch tasks)';
      }
    },

    async addTodo() {
      if (this.newTask.trim() === '') return;

      try {
        const response = await fetch('/api/v1/todos', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: this.newTask,
            Status: 'created',
            Priority: this.newPriority,
          }),
        });

        if (!response.ok) throw new Error(`Failed to create todo. statusCode: ${response.status}`);

        const data = await response.json();
        this.todos.push(data.data);
        this.newTask = '';
        this.newPriority = 'low'; // Reset priority to default
        this.fetchTodos()
        this.setStatusMessage('タスクが追加されました (Task added)');
      } catch (error) {
        console.error('Error creating todo:', error);
        this.setStatusMessage('タスクの作成に失敗しました (Failed to create task)');
      }
    },

    // Method to set status message for a limited time
    setStatusMessage(message) {
      this.statusMessage = message;
      setTimeout(() => {
        this.statusMessage = '';
      }, 5000); // Clear message after 10 seconds
    },
    async editTodo(todo) {
      todo.isEditing = false;

      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: todo.Task,
            Priority: todo.Priority,
          }),
        });

        if (!response.ok) throw new Error(`Failed to edit todo. statusCode: ${response.status}`);
        this.setStatusMessage('タスクが編集されました (Task edited)');
      } catch (error) {
        console.error('Error editing todo:', error);
        this.setStatusMessage('タスクの編集に失敗しました (Failed to edit task)');
      }
    },

    enableEdit(todo) {
      // Set the editing state to true for the selected todo
      todo.isEditing = true;
    },

    async deleteTodo(id) {
      try {
        const response = await fetch(`/api/v1/todos/${id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) throw new Error(`Failed to delete todo. statusCode: ${response.status}`);

        this.todos = this.todos.filter((todo) => todo.ID !== id);
        this.setStatusMessage('タスクが削除されました (Task deleted)');
        this.filterTodos();
      } catch (error) {
        console.error('Error deleting todo:', error);
        this.setStatusMessage('タスクの削除に失敗しました (Failed to delete task)');
      }
    },
    sortTodos() {
      const priorityOrder = { low: 3, medium: 2, high: 1 }; // Define priority order
      this.todos.sort((a, b) => priorityOrder[a.Priority] - priorityOrder[b.Priority]); // Sort based on the priority string
    },
    async updateStatus(todo) {
      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            Status: todo.Status === 'done' ? 'created' : 'done',
          }),
        });

        if (!response.ok) throw new Error(`Failed to update todo status. statusCode: ${response.status}`);

        todo.Status = todo.Status === 'done' ? 'created' : 'done';
        this.fetchTodos()
        this.setStatusMessage('タスクのステータスが変更されました (Task status updated)', 5000); // Set message for 5 seconds
        this.filterTodos(); // Re-filter after updating status
      } catch (error) {
        console.error('Error updating todo status:', error);
        this.setStatusMessage('ステータスの更新に失敗しました (Failed to update status)', 5000); // Set error message for 5 seconds
      }
    },

    // Update the setStatusMessage method to accept duration
    setStatusMessage(message, duration = 10000) { // Default duration to 10 seconds
      this.statusMessage = message;
      setTimeout(() => {
        this.statusMessage = '';
      }, duration); // Clear message after specified duration
    },

    async deleteTodo(id) {
      try {
        const response = await fetch(`/api/v1/todos/${id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) throw new Error(`Failed to delete todo. statusCode: ${response.status}`);

        this.todos = this.todos.filter((todo) => todo.ID !== id);
        this.statusMessage = 'タスクが削除されました (Task deleted)';
        this.filterTodos(); // Re-filter after deletion
      } catch (error) {
        console.error('Error deleting todo:', error);
        this.statusMessage = 'タスクの削除に失敗しました (Failed to delete task)';
      }
    },
    filterTodos() {
      // Filter todos based on search criteria
      this.filteredTodos = this.todos.filter((todo) => {
        const matchesTask = todo.Task.toLowerCase().includes(this.searchTask.toLowerCase());
        const matchesStatus = this.searchStatus ? todo.Status === this.searchStatus : true;
        return matchesTask && matchesStatus;
      });
    },
    resetSearchInputs() {
      // Clear search inputs when changing search type or selecting a new status
      this.searchTask = '';
      this.searchStatus = '';
      this.filteredTodos = this.todos; // Reset to show all todos
    },
  },
};
</script>

<style scoped>
.todo-main {
  max-width: 500px; /* Increased width from 400px to 600px */
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-right: 10px;
}

button {
  padding: 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.green-button {
  background-color: #28a745; /* Green for Add and Search buttons */
  color: #fff;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  position: relative;
}

.priority-badge {
  position: absolute;
  left: 5px;
  top: 5px;
  padding: 3px 7px;
  border-radius: 4px;
  font-size: 12px;
  color: #fff;
  background-color: #007bff;
}

.priority-high {
  background-color: #dc3545; /* Red for high priority */
}

.priority-medium {
  background-color: #ffc107; /* Yellow for medium priority */
}

.priority-low {
  background-color: #28a745; /* Green for low priority */
}

.buttons button {
  background-color: #f0f0f0;
  color: #333;
  margin-left: 5px;
  border-radius: 4px;
  padding: 5px 10px;
  transition: background-color 0.3s, color 0.3s;
}

.buttons button.done {
  background-color: #28a745; /* Green for done button */
  color: #fff;
}

.buttons button.delete-button {
  background-color: #f08080; /* Light Coral for delete button */
  color: white;
}

.status-message {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f0f0f0;
  border-radius: 4px;
}

.done-task {
  text-decoration: line-through;
  color: #6c757d;
}

.edit-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.search-section {
  display: flex;
  align-items: stretch; /* Ensure all elements stretch to the same height */
  gap: 10px; /* Gap between the elements */
}

.search-dropdown,
.search-input,
.search-dropdown {
  width: 150px; /* Set a fixed width to ensure both have the same size */
  padding: 8px; /* Add padding to make it visually consistent */
  border: 1px solid #ccc; /* Add border for consistency */
  border-radius: 4px; /* Optional: To add a consistent border radius */
  font-size: 14px; /* Set a consistent font size */
  box-sizing: border-box; /* Ensure padding doesn't affect the overall size */
}


/* Specific styles for button */
.green-button {
  background-color: #28a745; /* Green for Search button */
  color: #fff;
  border: none; /* Remove border */
  cursor: pointer; /* Pointer cursor for button */
  height: 40px; /* Ensure button height matches input fields */
}

/* Optional: Add hover effect for the button */
.green-button:hover {
  background-color: #218838; /* Darker green on hover */
}

.search-input {
  flex: 1; /* Allow input to take available space */
  /* Remove the border to unify the style */
  border: 1px solid #ddd; /* Add border back for consistency */
}

.edit-dropdown {
  display: flex;
  justify-content: flex-end; /* Aligns content to the right */
}

.right-align {
  text-align: right; /* Aligns text */
  appearance: none; /* Remove default appearance */
  -webkit-appearance: none; /* Remove default appearance for WebKit */
  -moz-appearance: none; /* Remove default appearance for Firefox */
}

.status-dropdown {
  position: absolute; /* Keep it absolute if needed */
  right: 0; /* Align to the right */
  font-size: 16px; /* Increase text size */
  padding: 10px; /* Add padding for a larger box */
  width: 200px; /* Set a specific width for the dropdown */
  height: 40px; /* Set a specific height for the dropdown */
  border: 1px solid #ddd; /* Optional: Add border to the dropdown */
  border-radius: 4px; /* Optional: Rounded corners */
  background-color: white; /* Optional: Set background color */
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); /* Optional: Add shadow for depth */
}



.dropdown-container {
  display: flex;
  justify-content: flex-end; /* Aligns content to the right */
  align-items: center; /* Center items vertically */
  position: relative; /* Positioning for dropdowns */
  margin-left: 10px; /* Add space between the edit button and dropdown */
}

/* Additional styling to ensure dropdowns appear correctly */
.status-dropdown {
  position: absolute; /* Allow dropdown to be positioned */
  right: 0; /* Align dropdown to the right */
  top: 100%; /* Position below the parent */
  background: white; /* Ensures background visibility */
  z-index: 10; /* Ensure it appears above other content */
}


.dropdown-container select,
.dropdown-container input {
  margin-left: 5px; /* Space between the dropdown and input */
}


</style>




