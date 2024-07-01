package geo

import "fmt"

// Tri3d is shorthand for creating a triangle.
func Tri3d[T Number](a, b, c Point3d[T]) Triangle3d[T] {
	return Triangle3d[T]{A: a, B: b, C: c}
}

func Tri3dFlat[T Number](ax, ay, az, bx, by, bz, cx, cy, cz T) Triangle3d[T] {
	return Triangle3d[T]{A: Pt3d(ax, ay, az),
		B: Pt3d(bx, by, bz),
		C: Pt3d(cx, cy, cz)}
}

type Triangle3d[T Number] struct {
	A Point3d[T]
	B Point3d[T]
	C Point3d[T]
}

func (t Triangle3d[T]) Center() Point3d[T] {
	return Point3d[T]{X: (t.A.X + t.B.X + t.C.X) / 3.,
		Y: (t.A.Y + t.B.Y + t.C.Y) / 3.,
		Z: (t.A.Z + t.B.Z + t.C.Z) / 3.}
}

/*
import numpy as np

def calculate_triangle_normal(triangle):
  """
  Calculates the normal vector of a 3D triangle represented by its vertices.

  Args:
      triangle: A numpy array of shape (3, 3) representing the triangle vertices.

  Returns:
      A numpy array of shape (3,) representing the normal vector of the triangle.
  """

  # Get the triangle vertices
  v1, v2, v3 = triangle

  # Calculate two vectors along the edges of the triangle
  edge1 = v2 - v1
  edge2 = v3 - v1

  # Calculate the normal vector using the cross product
  normal = np.cross(edge1, edge2)

  # Normalize the normal vector (optional)
  normal = normal / np.linalg.norm(normal)

  return normal

# Example usage
triangle = np.array([[1, 0, 0], [2, 1, 0], [1, 2, 0]])
normal = calculate_triangle_normal(triangle)

print("Triangle normal:", normal)
*/

/*
	func (t Triangle3d[T]) Normal() Point3d[T] {
		U, V := t.B.Sub(t.A), t.C.Sub(t.A)
		x := (U.Y * V.Z) - (U.Z * V.Y)
		y := (U.Z * V.X) - (U.X * V.Z)
		z := (U.X * V.Y) - (U.Y * V.X)
		return Point3d[T]{X: x, Y: y, Z: z}
	}
*/
func (t Triangle3d[T]) Normal() Point3d[T] {
	U, V := t.B.Sub(t.A), t.C.Sub(t.A)
	return U.Cross(V)
	// return U.Cross(V).Normalize()
}

/*
func TriangleIntersectionPlaneZ(tri Tri3dF, planeZ float64) (Pt3dF, bool) {
	triangleNormal := tri.Normal()

	// Ensure the triangle normal has a positive z-component (optional for some applications)
	if triangleNormal.Z < 0 {
		triangleNormal.X *= -1
		triangleNormal.Y *= -1
		triangleNormal.Z *= -1
	}

	// Check if triangle is degenerate (all points colinear)
	if Pt3dNearZero(triangleNormal, 1e-6) {
		return Pt3dF{}, false
	}
}
*/

func TrianglePlaneIntersection(tri Tri3dF, triPt Pt3dF, planeTri Tri3dF, planePt Pt3dF) (Pt3dF, bool) {
	tri_normal, plane_normal := tri.Normal().Normalize(), planeTri.Normal()
	lineA, lineB := tri.A, tri.A.Add(tri_normal)
	return LinePlaneIntersection(lineA, lineB, plane_normal, planePt)
}

func LinePlaneIntersection(lineA, lineB Pt3dF, planeNormal, planeOrigin Pt3dF) (Pt3dF, bool) {
	//  Line line,
	//  Plane plane,
	//  out double lineParameter )
	//{
	//  XYZ planePoint = plane.Origin;
	//  XYZ planeNormal = plane.Normal;
	linePoint := lineA
	lineDirection := (lineB.Sub(lineA)).Normalize()
	// Is the line parallel to the plane, i.e.,
	// perpendicular to the plane normal?

	if Float64NearZero(planeNormal.Dot(lineDirection), 1e-6) {
		return Pt3dF{}, false
	}

	lineParameter := planeNormal.Dot(planeOrigin) - planeNormal.Dot(linePoint)/planeNormal.Dot(lineDirection)

	return linePoint.Add(lineDirection.MultT(lineParameter)), true
}

/*
public static XYZ LinePlaneIntersection(
  Line line,
  Plane plane,
  out double lineParameter )
{
  XYZ planePoint = plane.Origin;
  XYZ planeNormal = plane.Normal;
  XYZ linePoint = line.GetEndPoint( 0 );

  XYZ lineDirection = (line.GetEndPoint( 1 )
    - linePoint).Normalize();

  // Is the line parallel to the plane, i.e.,
  // perpendicular to the plane normal?

  if( IsZero( planeNormal.DotProduct( lineDirection ) ) )
  {
    lineParameter = double.NaN;
    return null;
  }

  lineParameter = (planeNormal.DotProduct( planePoint )
    - planeNormal.DotProduct( linePoint ))
      / planeNormal.DotProduct( lineDirection );

  return linePoint + lineParameter * lineDirection;
}
*/

//F intersection_point(ray_direction, ray_point, plane_normal, plane_point)
//   R ray_point - ray_direction * dot(ray_point - plane_point, plane_normal) / dot(ray_direction, plane_normal)

func TrianglePlaneIntersection2(tri Tri3dF, triPt Pt3dF, planeTri Tri3dF, planePt Pt3dF) (Pt3dF, bool) {
	tri_normal, plane_normal := tri.Normal(), planeTri.Normal()

	tri_normal = tri_normal.Normalize()
	plane_normal = plane_normal.Normalize()
	u := triPt.Sub(planePt).Dot(plane_normal)
	v := tri_normal.Dot(plane_normal)
	pt := tri_normal.MultT(u).DivT(v)
	return pt, true
	//   R ray_point - ray_direction * dot(ray_point - plane_point, plane_normal) / dot(ray_direction, plane_normal)

	//   R ray_point - ray_direction * dot(ray_point - plane_point, plane_normal) / dot(ray_direction, plane_normal)

}

// Proiject a ray perpendicular to triPt on tri to find
// its intersection with the plane (as defined by planeTri and planePt)
func TrianglePlaneIntersection5(tri Tri3dF, triPt Pt3dF, planeTri Tri3dF, planePt Pt3dF) (Pt3dF, bool) {
	//func TrianglePlaneIntersection(tri Tri3dF, triPt Pt3dF, plane_normal, planePt Pt3dF) (Pt3dF, bool) {
	//	def project_ray_to_plane(triangle, point_on_triangle, plane_normal, plane_point):

	tri_normal, plane_normal := tri.Normal(), planeTri.Normal()
	//	tri_normal := tri.Normal()
	// Check if triPt is coplanar with the triangle
	//	if tri_normal.DotProduct(triPt.Sub(tri.A)) == 0. {
	//		fmt.Println("triPt coplanar")
	//		return Pt3dF{}, false
	//	}

	fmt.Println("plane_normal", plane_normal)
	// Direction of the ray (perpendicular to triangle normal)
	ray_direction := tri_normal
	fmt.Println("ray_direction", ray_direction)
	// Vector from plane point to point on triangle
	point_diff := triPt.Sub(planePt)
	fmt.Println("point_diff", point_diff)

	// Project point_diff onto the plane normal
	//	projection = np.dot(point_diff, plane_normal) / np.dot(plane_normal, plane_normal) * plane_normal
	projection := plane_normal.MultT(point_diff.DotProduct(plane_normal) / plane_normal.DotProduct(plane_normal))
	fmt.Println("projection", projection, "pd.dot", point_diff.DotProduct(plane_normal), "planeDot", plane_normal.DotProduct(plane_normal))

	// Distance along ray direction to reach the projection
	t := projection.DotProduct(plane_normal) / ray_direction.DotProduct(plane_normal)
	fmt.Println("t", t, projection.DotProduct(plane_normal), ray_direction.DotProduct(plane_normal))

	// If t is negative, the intersection is behind the point on triangle
	if t < 0. {
		fmt.Println("no projection")
		return Pt3dF{}, false
	}

	// Intersection point
	intersection := triPt.Add(ray_direction.MultT(t))
	fmt.Println("intersection", intersection)

	return intersection, true
}

/*
import numpy as np

def project_ray_to_plane(triangle, point_on_triangle, plane_normal, plane_point):
  """
  Projects a ray perpendicular to a triangle from a point on the triangle
  and finds its intersection with a plane.

  Args:
      triangle: A numpy array of shape (3, 3) representing the triangle vertices.
      point_on_triangle: A numpy array of shape (3,) representing a point on the triangle.
      plane_normal: A numpy array of shape (3,) representing the normal vector of the plane.
      plane_point: A numpy array of shape (3,) representing a point on the plane.

  Returns:
      A numpy array of shape (3,) representing the intersection point or None if no intersection exists.
  """

  # Calculate triangle normal
  v1, v2, v3 = triangle
  triangle_normal = np.cross(v2 - v1, v3 - v1)

  # Check if point_on_triangle is coplanar with the triangle
  if np.dot(triangle_normal, point_on_triangle - v1) == 0:
    return None

  # Direction of the ray (perpendicular to triangle normal)
  ray_direction = triangle_normal

  # Vector from plane point to point on triangle
  point_diff = point_on_triangle - plane_point

  # Project point_diff onto the plane normal
  projection = np.dot(point_diff, plane_normal) / np.dot(plane_normal, plane_normal) * plane_normal

  # Distance along ray direction to reach the projection
  t = np.dot(projection, plane_normal) / np.dot(ray_direction, plane_normal)

  # If t is negative, the intersection is behind the point on triangle
  if t < 0:
    return None

  # Intersection point
  intersection = point_on_triangle + t * ray_direction

  return intersection

# Example usage
triangle = np.array([[1, 0, 0], [2, 1, 0], [1, 2, 0]])
point_on_triangle = np.array([1.5, 0.5, 0])
plane_normal = np.array([0, 1, 1])
plane_point = np.array([0, 0, 0])

intersection = project_ray_to_plane(triangle, point_on_triangle, plane_normal, plane_point)

if intersection is not None:
  print("Intersection point:", intersection)
else:
  print("No intersection found")

*/

type Tri3dF = Triangle3d[float64]
