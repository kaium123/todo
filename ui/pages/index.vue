<template>
  <div class="todo-main">
    <h1>TODOリスト</h1>

    <!-- Show status messages -->
    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>

    <!-- Input group to add tasks -->
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

    <!-- Search Section -->
    <div class="search-section">
      <select v-model="searchType" class="search-dropdown" @change="resetSearchInputs">
        <option value="task">タスク (Task)</option>
        <option value="status">ステータス (Status)</option>
      </select>
      <input
        v-if="searchType === 'task'"
        v-model="searchTask"
        placeholder="タスクを検索 (Search Task)"
        class="search-input"
        @input="filterTodos"
      />
      <select
        v-if="searchType === 'status'"
        v-model="searchStatus"
        class="search-dropdown"
        @change="filterTodos"
      >
        <option value="">すべて (All)</option>
        <option value="created">作成済み (Created)</option>
        <option value="processing">処理中 (Processing)</option>
        <option value="done">完了 (Done)</option>
      </select>
      <button
        class="green-button"
        @click="searchType === 'task' ? searchByTask() : searchByStatus()"
      >
        検索 (Search)
      </button>
    </div>

    <!-- Display tasks -->
    <div v-if="filteredTodos.length > 0">
      <div v-for="todo in filteredTodos" :key="todo.ID" class="todo-item">
        <div>
          <span
            v-if="!todo.isEditing"
            :class="{ 'done-task': todo.Status === 'done' }"
            @click="updateStatus(todo)"
          >
            {{ todo.Task }} {{ todo.Status }}
          </span>
          <div v-else class="edit-dropdown">
            <input
              v-model="todo.Task"
              class="edit-input"
              @blur="editTodo(todo)"
              @keyup.enter="editTodo(todo)"
              placeholder="Edit Task"
            />
            <select
              v-model="todo.Status"
              class="status-dropdown"
              @change="editTodo(todo)"
            >
              <option value="created">作成済み (Created)</option>
              <option value="processing">処理中 (Processing)</option>
              <option value="done">完了 (Done)</option>
            </select>
          </div>
        </div>
        <div class="buttons">
          <button class="edit-button" @click="enableEdit(todo)">✏️</button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
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
        this.fetchTodos();
        this.setStatusMessage('タスクが追加されました (Task added)');
      } catch (error) {
        console.error('Error creating todo:', error);
        this.setStatusMessage('タスクの作成に失敗しました (Failed to create task)');
      }
    },

    setStatusMessage(message) {
      this.statusMessage = message;
      setTimeout(() => {
        this.statusMessage = '';
      }, 5000); // Clear message after 5 seconds
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
            Status: todo.Status, // Include status in the update
            Priority: todo.Priority, // You may need to send the priority if it changes
          }),
        });

        this.fetchTodos()
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
        this.fetchTodos();
        this.setStatusMessage('タスクのステータスが変更されました (Task status updated)'); // Set message for 5 seconds
        this.filterTodos(); // Re-filter after updating status
      } catch (error) {
        console.error('Error updating todo status:', error);
        this.setStatusMessage('ステータスの更新に失敗しました (Failed to update status)'); // Set error message for 5 seconds
      }
    },
  },
};
</script>

<style scoped>
.todo-main {
  max-width: 500px; /* Increased width from 400px to 600px */
  margin: 20px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 8px;
}

.input-group,
.search-section {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.input-group input,
.input-group select,
.search-section select,
.search-section input {
  flex: 1; /* Allow these elements to grow */
  padding: 10px;
  margin-right: 10px;
}

.green-button {
  padding: 10px 15px;
  background-color: #28a745;
  border: none;
  color: white;
  border-radius: 5px;
  cursor: pointer;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.done-task {
  text-decoration: line-through; /* Strike-through for completed tasks */
}

.edit-dropdown {
  display: flex;
  align-items: center;
}

.edit-input,
.status-dropdown {
  margin-right: 10px;
}

.buttons {
  display: flex;
  gap: 10px;
}

.edit-button,
.delete-button {
  background: none;
  border: none;
  cursor: pointer;
  color: #007bff;
}

.edit-button:hover,
.delete-button:hover {
  color: #0056b3;
}

.status-message {
  color: green;
  margin-bottom: 10px;
}
</style>
