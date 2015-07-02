package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type TestRecord struct {
	Title string `csv:"title"`
}

type TestRecordWithoutTag struct {
	Title string
}

func TestGetFieldIndexFromHeader(t *testing.T) {
	result, err := getFieldIndexFromHeader(TestRecord{}, []string{"title"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result[0])

	result, err = getFieldIndexFromHeader(TestRecord{}, []string{"another", "title"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result[0])

	result, err = getFieldIndexFromHeader(TestRecord{}, []string{"another", "Title"})
	assert.Error(t, err, "Should throw error if no match header")

	result, err = getFieldIndexFromHeader(TestRecordWithoutTag{}, []string{"Another", "Title"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result[0], "Should use field name")
}

type TestFillRecord struct {
	Title string `csv:"title"`
	Count int    `csv:"count"`
}

func TestFillStruct(t *testing.T) {
	s := TestRecord{}
	fieldColumnIndex, _ := getFieldIndexFromHeader(s, []string{"another", "title"})

	err := fillStruct(&s, fieldColumnIndex, []string{"test", "value"})
	assert.NoError(t, err)
	assert.Equal(t, "value", s.Title)
}

func TestFillStruct_Int(t *testing.T) {
	s := TestFillRecord{}
	fieldColumnIndex, _ := getFieldIndexFromHeader(s, []string{"another", "count", "title"})

	err := fillStruct(&s, fieldColumnIndex, []string{"test", "2", "value"})
	assert.NoError(t, err)
	assert.Equal(t, "value", s.Title)
	assert.Equal(t, 2, s.Count)

	err = fillStruct(&s, fieldColumnIndex, []string{"test", "value", "2"})
	assert.Error(t, err, "Should throw error if data type not match with field type")
}

type TestFillRecordUnSupportedType struct {
	Value float32
}

func TestFillStruct_NotImplementedType(t *testing.T) {
	s := TestFillRecordUnSupportedType{}
	fieldColumnIndex, _ := getFieldIndexFromHeader(s, []string{"Value"})

	err := fillStruct(&s, fieldColumnIndex, []string{"3.14"})
	assert.Error(t, err)
}

type TestFillRecordTime struct {
	Time time.Time `format:"2006-01-02 03:04:05"`
}

func TestFillStruct_Time(t *testing.T) {
	s := TestFillRecordTime{}
	fieldColumnIndex, _ := getFieldIndexFromHeader(s, []string{"Time"})

	err := fillStruct(&s, fieldColumnIndex, []string{"2015-06-10 09:21:43"})
	assert.NoError(t, err)
	assert.EqualValues(t, 6, s.Time.Month())
	assert.Equal(t, 10, s.Time.Day())
	assert.Equal(t, 2015, s.Time.Year())
	assert.Equal(t, 9, s.Time.Hour())
	assert.Equal(t, 21, s.Time.Minute())
	assert.Equal(t, 43, s.Time.Second())
}

func TestReadReport(t *testing.T) {
	fixture := `错误摘要	发生次数	首次发生时间	版本	错误详情
java.lang.NullPointerException	8	2015-06-10 09:21:43	3.0.3	"java.lang.NullPointerException
	at android.support.v7.widget.RecyclerView.onInterceptTouchEvent(Unknown Source)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1859)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at android.view.ViewGroup.dispatchTransformedTouchEvent(ViewGroup.java:2218)
	at android.view.ViewGroup.dispatchTouchEvent(ViewGroup.java:1917)
	at com.android.internal.policy.impl.PhoneWindow$DecorView.superDispatchTouchEvent(PhoneWindow.java:2117)
	at com.android.internal.policy.impl.PhoneWindow.superDispatchTouchEvent(PhoneWindow.java:1564)
	at android.app.Activity.dispatchTouchEvent(Activity.java:2493)
	at android.support.v7.internal.view.WindowCallbackWrapper.dispatchTouchEvent(Unknown Source)
	at android.support.v7.internal.view.WindowCallbackWrapper.dispatchTouchEvent(Unknown Source)
	at com.android.internal.policy.impl.PhoneWindow$DecorView.dispatchTouchEvent(PhoneWindow.java:2065)
	at android.view.View.dispatchPointerEvent(View.java:7903)
	at android.view.ViewRootImpl$ViewPostImeInputStage.processPointerEvent(ViewRootImpl.java:4188)
	at android.view.ViewRootImpl$ViewPostImeInputStage.onProcess(ViewRootImpl.java:4067)
	at android.view.ViewRootImpl$InputStage.deliver(ViewRootImpl.java:3624)
	at android.view.ViewRootImpl$InputStage.onDeliverToNext(ViewRootImpl.java:3674)
	at android.view.ViewRootImpl$InputStage.forward(ViewRootImpl.java:3643)
	at android.view.ViewRootImpl$AsyncInputStage.forward(ViewRootImpl.java:3750)
	at android.view.ViewRootImpl$InputStage.apply(ViewRootImpl.java:3651)
	at android.view.ViewRootImpl$AsyncInputStage.apply(ViewRootImpl.java:3807)
	at android.view.ViewRootImpl$InputStage.deliver(ViewRootImpl.java:3624)
	at android.view.ViewRootImpl$InputStage.onDeliverToNext(ViewRootImpl.java:3674)
	at android.view.ViewRootImpl$InputStage.forward(ViewRootImpl.java:3643)
	at android.view.ViewRootImpl$InputStage.apply(ViewRootImpl.java:3651)
	at android.view.ViewRootImpl$InputStage.deliver(ViewRootImpl.java:3624)
	at android.view.ViewRootImpl.deliverInputEvent(ViewRootImpl.java:5836)
	at android.view.ViewRootImpl.doProcessInputEvents(ViewRootImpl.java:5816)
	at android.view.ViewRootImpl.enqueueInputEvent(ViewRootImpl.java:5787)
	at android.view.ViewRootImpl$WindowInputEventReceiver.onInputEvent(ViewRootImpl.java:5925)
	at android.view.InputEventReceiver.dispatchInputEvent(InputEventReceiver.java:185)
	at android.os.MessageQueue.nativePollOnce(Native Method)
	at android.os.MessageQueue.next(MessageQueue.java:138)
	at android.os.Looper.loop(Looper.java:123)
	at android.app.ActivityThread.main(ActivityThread.java:5314)
	at java.lang.reflect.Method.invokeNative(Native Method)
	at java.lang.reflect.Method.invoke(Method.java:515)
	at com.android.internal.os.ZygoteInit$MethodAndArgsCaller.run(ZygoteInit.java:864)
	at com.android.internal.os.ZygoteInit.main(ZygoteInit.java:680)
	at dalvik.system.NativeStart.main(Native Method)
"`
	results, err := readReport(strings.NewReader(fixture))
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 8, results[0].Count)
}

func TestReadFile(t *testing.T) {
	file, err := ReadFile(filepath.Join("fixtures", "test.csv"))
	assert.NoError(t, err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	lines := strings.Split(string(content), "\n")
	header := strings.Split(lines[0], "\t")
	assert.Equal(t, 5, len(header))
	assert.Equal(t, "错误摘要", header[0])
	assert.Equal(t, "发生次数", header[1])
	fieldColumnIndex, err := getFieldIndexFromHeader(Record{}, header)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(fieldColumnIndex))
	assert.Equal(t, 0, fieldColumnIndex[0])
	assert.Equal(t, 1, fieldColumnIndex[1])
	assert.Equal(t, 2, fieldColumnIndex[2])
	assert.Equal(t, 3, fieldColumnIndex[3])
	assert.Equal(t, 4, fieldColumnIndex[4])

	file, err = ReadFile(filepath.Join("fixtures", "test.csv"))
	assert.NoError(t, err)

	results, err := readReport(file)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 1, len(results))
	assert.Equal(t, 8, results[0].Count)
}
