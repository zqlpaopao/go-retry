package pkg

type TaskFunc func(RetryableFunc, ...interface{})

//Task 任务元数据
type Task struct {
	f            TaskFunc
	retryableFun RetryableFunc
	CustomParams []interface{}
}

//InitTask -- ----------------------------
//--> @Description 初始化任务
//--> @Param
//--> @return
//-- ----------------------------
//func InitTask(argF *Task)*Task {
//	return &Task{
//		f: argF,
//		retryableFun:
//	}
//}

//Execute -- ----------------------------
//--> @Description 任务执行主体
//--> @Param
//--> @return
//-- ----------------------------
func (t *Task) Execute() {
	t.f(t.retryableFun, t.CustomParams...)
}

/*****************************************协程池角色*******************************/
type pool struct {
	receiveCh chan *Task
	runCh     chan *Task
	workerNum int
}

//NewPool -- ----------------------------
//--> @Description 初始化协程池
//--> @Param
//--> @return
//-- ----------------------------
func NewPool(n int) *pool {
	return &pool{
		receiveCh: make(chan *Task),
		runCh:     make(chan *Task),
		workerNum: n,
	}
}

func (p *pool) AddTask(f *Task) {
	p.receiveCh <- f
}

//-- ----------------------------
//--> @Description worker执行器
//--> @Param
//--> @return
//-- ----------------------------
func (p *pool) worker(i int) {
	for task := range p.runCh {
		task.Execute()
	}
}

//Run -- ----------------------------
//--> @Description 协程工作初始化
//--> @Param
//--> @return
//-- ----------------------------
func (p *pool) Run() {
	//启动指定数量的协程
	for i := 0; i < p.workerNum; i++ {
		go p.worker(i)
	}

	//将接收的任务同步给内部runChan
	for task := range p.receiveCh {
		p.runCh <- task
	}
}
