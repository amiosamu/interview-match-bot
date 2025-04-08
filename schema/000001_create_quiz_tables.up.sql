-- Create quiz tables
CREATE TABLE IF NOT EXISTS quiz_questions (
    id SERIAL PRIMARY KEY,
    language VARCHAR(50) NOT NULL,
    category VARCHAR(50) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    question_text TEXT NOT NULL,
    answer_options JSONB NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    explanation TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User quiz sessions table
CREATE TABLE IF NOT EXISTS user_quiz_sessions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    language VARCHAR(50) NOT NULL,
    current_question_index INT NOT NULL DEFAULT 0,
    question_ids JSONB NOT NULL,
    correct_answers INT NOT NULL DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    UNIQUE (user_id, language, started_at)
);

-- User quiz answers table
CREATE TABLE IF NOT EXISTS user_quiz_answers (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    session_id INT NOT NULL REFERENCES user_quiz_sessions(id),
    question_id INT NOT NULL REFERENCES quiz_questions(id),
    answer_given VARCHAR(255) NOT NULL,
    is_correct BOOLEAN NOT NULL,
    answered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_quiz_questions_language ON quiz_questions(language);
CREATE INDEX IF NOT EXISTS idx_user_quiz_sessions_user_id ON user_quiz_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_quiz_sessions_active ON user_quiz_sessions(user_id) WHERE completed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_quiz_answers_session ON user_quiz_answers(session_id);

-- Seed Golang Questions
INSERT INTO quiz_questions (language, category, difficulty, question_text, answer_options, correct_answer, explanation)
VALUES
-- Basic syntax questions
('golang', 'syntax', 'beginner', 
 'What is the zero value of an integer type in Go?', 
 '["0", "nil", "undefined", "null"]', 
 '0', 
 'In Go, all variables are initialized to their zero value. For numeric types like int, float, etc., the zero value is 0.'),

('golang', 'syntax', 'beginner', 
 'Which of the following is a valid variable declaration in Go?', 
 '["var x int = 10", "let x = 10", "x := 10", "Both A and C"]', 
 'Both A and C', 
 'In Go, you can declare variables using the var keyword with an explicit type (var x int = 10) or using the short declaration operator := which infers the type (x := 10).'),

('golang', 'syntax', 'intermediate', 
 'What does the following code print?\n\nfunc main() {\n  x := 10\n  {\n    x := 20\n    fmt.Println(x)\n  }\n  fmt.Println(x)\n}', 
 '["10\\n10", "20\\n20", "20\\n10", "10\\n20"]', 
 '20\n10', 
 'This demonstrates variable shadowing in Go. Inside the inner block, a new variable x is declared which shadows the outer x. After the inner block ends, the outer x is visible again.'),

-- Functions and methods questions
('golang', 'functions', 'beginner', 
 'What is the correct way to define a function that returns an integer in Go?', 
 '["function sum(a, b int) int { return a + b }", "func sum(a int, b int) -> int { return a + b }", "func sum(a, b int) int { return a + b }", "def sum(a, b int) int { return a + b }"]', 
 'func sum(a, b int) int { return a + b }', 
 'In Go, functions are defined using the func keyword. Parameters of the same type can be grouped (a, b int), and the return type comes after the parameter list.'),

('golang', 'functions', 'intermediate', 
 'What is a defer statement used for in Go?', 
 '["To handle errors", "To delay execution of a function until the surrounding function returns", "To create goroutines", "To define interfaces"]', 
 'To delay execution of a function until the surrounding function returns', 
 'The defer statement pushes a function call onto a list. The list of saved calls is executed after the surrounding function returns. Defers are commonly used for cleanup operations.'),

-- Concurrency questions
('golang', 'concurrency', 'intermediate', 
 'What does the following code do?\n\ngo func() {\n  fmt.Println("Hello")\n}()', 
 '["Creates a new thread", "Executes the function synchronously", "Creates a new goroutine", "Causes a compilation error"]', 
 'Creates a new goroutine', 
 'The go keyword before a function call creates a new goroutine, which is a lightweight thread managed by the Go runtime. The function executes concurrently with the calling function.'),

('golang', 'concurrency', 'advanced', 
 'What is the primary purpose of channels in Go?', 
 '["To allocate memory", "To synchronize goroutines and enable communication between them", "To handle errors across functions", "To define interfaces for types"]', 
 'To synchronize goroutines and enable communication between them', 
 'Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values in another goroutine, providing both communication and synchronization.'),

-- Error handling questions
('golang', 'error', 'intermediate', 
 'What is the idiomatic way to handle errors in Go?', 
 '["Try-catch blocks", "Return error values and check them", "Exception handling", "Using panic and recover"]', 
 'Return error values and check them', 
 'Go does not have exceptions. Instead, functions return an error value that the caller should check. This explicit error handling is a core philosophy of Go programming.'),

('golang', 'error', 'intermediate', 
 'Which statement about panic in Go is correct?', 
 '["Panic is Gos version of exceptions", "Panic should be used for routine error handling", "Panic causes the program to exit immediately", "Panic unwinds the stack, executing deferred functions"]', 
 'Panic unwinds the stack, executing deferred functions', 
 'When a function panics, normal execution stops, deferred functions are executed, and control returns to the caller. This continues up the stack until all functions in the goroutine have returned, at which point the program crashes.'),

-- Interfaces and types questions
('golang', 'types', 'intermediate', 
 'How does a type satisfy an interface in Go?', 
 '["By explicitly declaring that it implements the interface", "By implementing all the methods required by the interface", "By extending the interface", "By using the implements keyword"]', 
 'By implementing all the methods required by the interface', 
 'In Go, a type implements an interface by implementing its methods. There is no explicit declaration of intent, no "implements" keyword. This is known as structural typing or duck typing.');