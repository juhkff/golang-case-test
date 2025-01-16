struct Test
{
    int a;
    float b;
    double type;
    int size : 10;
    int arr1[10];
    int arr2[];
};
typedef struct Test Test;

int Test_arr2_helper(struct Test *tm, int pos);
int get_size(struct Test *tm);
void set_size(struct Test *tm, int value);
enum Color
{
    RED,
    GREEN,
    BLUE
};
typedef enum Color Color;
#pragma pack(1)
struct Test2
{
    float a;
    char b;
    int c;
};