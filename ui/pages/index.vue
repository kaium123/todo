<template>
  <div class="todo-main">
    <h1>TODOリスト / TODO List</h1>
    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>

    <!-- Input group for new tasks with priority selection -->
    <div class="input-group">
      <input v-model="newTask" placeholder="新しいタスクを入力 / Enter a new task" @keyup.enter="addTodo" />
      <select v-model="newPriority">
        <option value="1">高 (High)</option>
        <option value="2">中 (Medium)</option>
        <option value="3">低 (Low)</option>
      </select>
      <button @click="addTodo">追加 / Submit</button> <!-- Submit button with both Japanese and English -->
    </div>

    <!-- Incomplete Tasks Section -->
    <div v-if="incompleteTodos.length > 0">
      <h2>未完了のタスク / Incomplete Tasks</h2>
      <div v-for="todo in sortedIncompleteTodos" :key="todo.ID" class="todo-item">
        <input
          v-if="todo.isEditing"
          v-model="todo.Task"
          class="edit-input"
          @blur="editTodo(todo)"
          @keyup.enter="editTodo(todo)"
        />
        <span v-else :class="{ 'done-task': todo.Status === 'done' }" @click="enableEdit(todo)">
          {{ todo.Task }} (優先度: {{ formatPriority(todo.Priority) }})
        </span>
        <div class="buttons">
          <button :class="{ 'done': todo.Status === 'done' }" @click="updateStatus(todo)">✔️</button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
        </div>
      </div>
    </div>

    <!-- Completed Tasks Section -->
    <div v-if="completedTodos.length > 0">
      <h2>完了したタスク / Completed Tasks</h2>
      <div v-for="todo in completedTodos" :key="todo.ID" class="todo-item">
        <span class="done-task">{{ todo.Task }}</span>
        <div class="buttons">
          <button @click="updateStatus(todo)">↩️</button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
        </div>
      </div>
    </div>

    <!-- No Tasks Message -->
    <div v-else-if="todos.length === 0">
      <p>タスクがありません / No tasks available.</p>
    </div>
  </div>
</template>




<script>
export default {
  data() {
    return {
      newTask: '',
      newPriority: '2', // Default priority is Medium
      todos: [],
      statusMessage: '',
    };
  },
  mounted() {
    this.fetchTodos();
  },
  computed: {
    incompleteTodos() {
      return this.todos.filter(todo => todo.Status !== 'done');
    },
    completedTodos() {
      return this.todos.filter(todo => todo.Status === 'done');
    },
    sortedIncompleteTodos() {
      return this.incompleteTodos.sort((a, b) => a.Priority - b.Priority);
    },
  },
  methods: {
    async fetchTodos() {
      try {
        const response = await fetch(`/api/v1/todos`);
        if (!response.ok) throw new Error(`Failed to get todo list. statusCode: ${response.status}`);
        const data = await response.json();
        this.todos = data.data;
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
        this.newPriority = '2'; // Reset priority to default
        this.statusMessage = 'タスクが追加されました (Task added successfully)';
      } catch (error) {
        console.error('Error creating todo:', error);
        this.statusMessage = 'タスクの作成に失敗しました (Failed to create task)';
      }
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
        this.statusMessage = 'タスクが編集されました (Task edited successfully)';
      } catch (error) {
        console.error('Error editing todo:', error);
        this.statusMessage = 'タスクの編集に失敗しました (Failed to edit task)';
      }
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
        this.statusMessage = 'タスクのステータスが変更されました (Task status updated)';
      } catch (error) {
        console.error('Error updating todo status:', error);
        this.statusMessage = 'ステータスの更新に失敗しました (Failed to update status)';
      }
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
        this.todos = this.todos.filter(todo => todo.ID !== id);
        this.statusMessage = 'タスクが削除されました (Task deleted successfully)';
      } catch (error) {
        console.error('Error deleting todo:', error);
        this.statusMessage = 'タスクの削除に失敗しました (Failed to delete task)';
      }
    },
    formatPriority(priority) {
      switch (priority) {
        case '1':
          return '高 (High)';
        case '2':
          return '中 (Medium)';
        case '3':
          return '低 (Low)';
        default:
          return '未設定 (Not set)';
      }
    },
  },
};

</script>

<style scoped>
.todo-main {
  max-width: 600px; /* Increased width */
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
  background-color: #333;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
  flex-shrink: 0; /* Prevent button from shrinking */
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
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
  background-color: #333;
  color: #fff;
}

.buttons button.delete-button {
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
}

.edit-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
